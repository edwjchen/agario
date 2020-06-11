package router

//client in python <---> player server <---> region server

import (
	"encoding/binary"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"hash/fnv"
	"log"
	. "peer_to_peer/server/region_pb"
	"sort"
	"sync"
	"time"
)

type RegionChangeInfo struct {
	Successor bool 
	Join      bool 
	Prev      uint32 
	Curr      uint32
	PrevConn  *grpc.ClientConn 
	CurrConn  *grpc.ClientConn

}

type Router struct {
	lock sync.Mutex 
	playerLock sync.RWMutex 
	conns map[string] *grpc.ClientConn // map from ip of ip to connections
	playerConns map[string] *grpc.ClientConn
	iphash  map[string]uint32
	haship  map[uint32]string
	liveBacks []uint32 // stores hashes of ip
	Hash      uint32
	CurrPredecessor uint32
	CurrSuccessor   uint32
	RegionChange    chan RegionChangeInfo
	Ready     bool
}

type heartbeatOutput struct {
	ip   string
	conn *grpc.ClientConn
}

// Need to store hash of self
func (r *Router) Init(servers []string, ownAddr string) {
	r.haship = make(map[uint32]string)
	r.iphash = make(map[string]uint32)
	r.conns = make(map[string] *grpc.ClientConn)
	r.playerConns = make(map[string] *grpc.ClientConn)
	for _, ip := range(servers) {
		hasher := fnv.New32a()
		hasher.Write([]byte(ip))
		hash := uint32(hasher.Sum32())
		r.iphash[ip] = hash 
		r.haship[hash] = ip
		r.conns[ip] = nil
		if ip == ownAddr {
			r.Hash = hash
		}
	}
}

// this heartbeat function runs as a bg goroutine
func (r *Router) Heartbeat() {
	log.Println("Beating my heart")
	// log.Println(r.conns)
	retChan := make(chan heartbeatOutput, len(r.conns))
	ping := func(ip string, cxn *grpc.ClientConn) {
		if cxn == nil {
			conn, err := grpc.Dial(ip, grpc.WithInsecure())
			if err != nil {
				log.Println("Failed dail", err)
				retChan <- heartbeatOutput{
					ip:   ip,
					conn: nil,
				}
				return
			} else {
				cxn = conn
			}
		}
		// TODO
		client := NewRegionClient(cxn)
		_, err := client.Ping(context.Background(), &EmptyRequest{})
		if err != nil {
			// log.Println("Failed ping", err)
			retChan <- heartbeatOutput{
				ip:   ip,
				conn: nil,
			}
			return 
		}

		// log.Println("OK")
		retChan <- heartbeatOutput {
			ip:   ip,
			conn: cxn,
		}
	}

	for {
		<-time.Tick(time.Second)
		r.lock.Lock()
		for h, cxn := range r.conns {
			// send rpc
			// log.Println(h, cxn)
			go ping(h, cxn)
		}
		r.liveBacks = []uint32{}
		for i := 0; i < len(r.conns); i++ {
			status := <-retChan
			r.conns[status.ip] = status.conn
			if status.conn != nil {
				r.liveBacks = append(r.liveBacks, r.iphash[status.ip])
			}			
		}

		sort.Slice(r.liveBacks, func(i, j int) bool {
			return r.liveBacks[i] < r.liveBacks[j]
		})

		if r.Ready {
			newSuccessor := r.successor(r.Hash+1)
			newPredecessor := r.predecessor(r.Hash)

			if newSuccessor != r.CurrSuccessor {
				// handle one node case
				if newSuccessor == r.Hash {
					r.RegionChange <- RegionChangeInfo{
						Successor: true,
						Join:      true,
						PrevConn:  nil,
						CurrConn:  nil,
					}
				} else {
					r.OnSccessorChange(r.CurrSuccessor, newSuccessor)
				}
				r.CurrSuccessor = newSuccessor
			}

			if newPredecessor != r.CurrPredecessor {
				if newPredecessor != r.Hash {
					r.onPredecessorChange(r.CurrPredecessor, newPredecessor)
				}
				r.CurrPredecessor = newPredecessor
			}
		}
		r.lock.Unlock()
	}
}

func (r *Router) UpdateRingPos() {
	r.lock.Lock()
	r.CurrSuccessor = r.successor(r.Hash+1)
	r.CurrPredecessor = r.predecessor(r.Hash)
	r.Ready = true
	r.lock.Unlock()
}

// Returns GRPC connection
func (r *Router) Get(key uint32) (*grpc.ClientConn, *grpc.ClientConn) {
	// return grpc connection of head of chain
	hasher := fnv.New32a()
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, key)
	hasher.Write(b)
	hash := hasher.Sum32()

	primaryHash := r.Successor(hash)
	bkupHash := r.Successor(primaryHash + 1)
	//log.Println("Get:",key, primaryHash, bkupHash)
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.conns[r.haship[primaryHash]], r.conns[r.haship[bkupHash]]
}

func (r *Router) GetPlayerConn(addr string) *grpc.ClientConn {
	r.playerLock.RLock()
	cxn, ok := r.playerConns[addr]
	r.playerLock.RUnlock()

	if !ok {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			log.Println("Failed dail", err)
			return nil
		} else {
			cxn = conn
			r.playerLock.Lock()
			r.playerConns[addr] = conn
			r.playerLock.Unlock()
		}
	}
	return cxn
}

// precondition: clients of this lib will only attempt to invalidate existing client connection
func (r *Router) InvalidatePlayerConn(addr string) {
	r.playerLock.Lock()
	defer r.playerLock.Unlock()
	delete(r.playerConns, addr)
}
// Returns GRPC connection
func (r *Router) GetSuccessor() *grpc.ClientConn {
	// get whatever is after us in aliveBacks
	successorHash := r.Successor(r.Hash+1)
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.conns[r.haship[successorHash]]
}

func (r *Router) GetPredecessor() *grpc.ClientConn {
	// get whatever is after us in aliveBacks
	predecessorHash := r.Predecessor(r.Hash)
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.conns[r.haship[predecessorHash]]
}

func (r *Router) successor(h uint32) uint32 {

	if len(r.liveBacks) == 1 {
		return r.liveBacks[0]
	}

	for i := 0; i < len(r.liveBacks) - 1; i++ {
		if r.liveBacks[i] < h && h <= r.liveBacks[i+1] {
			return r.liveBacks[i+1]
		}
	}

	//wraparound
	return r.liveBacks[0]
}

func (r *Router) Successor(h uint32) uint32 {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.successor(h)
}

func (r *Router) predecessor(h uint32) uint32 {
	if len(r.liveBacks) == 1 {
		return r.liveBacks[0]
	}

	for i := 0; i < len(r.liveBacks) - 1; i++ {
		if r.liveBacks[i] < h && h <= r.liveBacks[i+1] {
			return r.liveBacks[i]
		}
	}

	//wraparound
	return r.liveBacks[len(r.liveBacks) - 1]
}

func (r *Router) Predecessor(h uint32) uint32 {
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.predecessor(h)
}

func (r *Router) OnSccessorChange(oldSucc, newSucc uint32) {

	var newdist uint32 = newSucc - r.Hash
	var olddist uint32 = oldSucc - r.Hash
	log.Println("OnsuccessorChange",r.Hash, oldSucc, newSucc)
	if newdist < olddist {
		r.RegionChange <- RegionChangeInfo{
			Successor: true,
			Join:      true,
			Prev:      oldSucc,
			Curr:      newSucc,
			PrevConn:  nil,
			CurrConn:  r.conns[r.haship[newSucc]],
		}
	} else {
		r.RegionChange <- RegionChangeInfo{
			Successor: true,
			Join:      false,
			Prev:      oldSucc,
			Curr:      newSucc, 
			PrevConn:  r.conns[r.haship[oldSucc]],
			CurrConn:  r.conns[r.haship[newSucc]],
		}
	}

}

func (r *Router) onPredecessorChange(oldPred, newPred uint32) {

	var newdist uint32 = r.Hash - newPred
	var olddist uint32 = r.Hash - oldPred
	log.Println("OnPredChange",r.Hash, oldPred, newPred)
	if newdist > olddist {
		r.RegionChange <- RegionChangeInfo{
			Successor: false,
			Join:      false, 
			Prev:      oldPred,
			Curr:      newPred, 
			PrevConn:  nil,
			CurrConn:  r.conns[r.haship[newPred]],
		}
	} else {
		r.RegionChange <- RegionChangeInfo{
			Successor: false,
			Join:      true, 
			Prev:      oldPred,
			Curr:      newPred,  
			PrevConn:  r.conns[r.haship[oldPred]],
			CurrConn:  r.conns[r.haship[newPred]],
		}
	}

}

//func before(seq1, seq2 uint32) bool {
//	return int32(seq1-seq2) < 0
//}

// func (r *Router) onSuccessorJoin() {
// 	r.lock.Lock()
// 	defer r.lock.Unlock()
// 	// move regions on curr node that hash to curr region to joined node
// }

// func (r *Router) onSuccessorLeave() {
// 	r.lock.Lock()
// 	defer r.lock.Unlock()
// 	// move regions on curr node that hash to curr region to new successor
// }

// func (r *Router) onPredecessorJoin() {
// 	r.lock.Lock()
// 	defer r.lock.Unlock()
// 	// remove regions on curr node that hash to prepredecessor
// 	// move regions that hash to joining node from curr node to joining node
// 	// remove regions on successor that hash to joining node
// }

// func (r *Router) onPredecessorLeave() {
// 	r.lock.Lock()
// 	defer r.lock.Unlock()
// 	// move regions that hashed to node that left to successor
// }