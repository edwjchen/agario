package router

//client in python <---> player server <---> region server

import (
	"hash/fnv"
	"sync"
	"log"
	"time"
	"encoding/binary"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"peer_to_peer/server/region_pb"
	"sort"
)

type Router struct {
	lock sync.Mutex 
	playerLock sync.RWMutex 
	conns map[string] *grpc.ClientConn // map from ip of ip to connections
	playerConns map[string] *grpc.ClientConn
	iphash map[string]uint32
	haship map[uint32]string
	liveBacks []uint32 // stores hashes of ip
	Hash uint32
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
		client := region_pb.NewRegionClient(cxn)
		_, err := client.Ping(context.Background(), &region_pb.EmptyRequest{})
		if err != nil {
			log.Println("Failed ping", err)
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

		r.lock.Unlock()
	}
}

// Returns GRPC connection
func (r *Router) Get(key uint32) *grpc.ClientConn {
	// return grpc connection of head of chain
	hasher := fnv.New32a()
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, key)
	hasher.Write(b)
	hash := uint32(hasher.Sum32())

	primaryHash := r.Successor(hash)
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.conns[r.haship[primaryHash]]
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

func (r *Router) Successor(h uint32) uint32 {
	r.lock.Lock()
	defer r.lock.Unlock()
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

func (r *Router) onSuccessorJoin() {
	r.lock.Lock()
	defer r.lock.Unlock()
	// move regions on curr node that hash to curr region to joined node
}

func (r *Router) onSuccessorLeave() {
	r.lock.Lock()
	defer r.lock.Unlock()
	// move regions on curr node that hash to curr region to new successor
}

func (r *Router) onPredecessorJoin() {
	r.lock.Lock()
	defer r.lock.Unlock()
	// remove regions on curr node that hash to prepredecessor
	// move regions that hash to joining node from curr node to joining node
	// remove regions on successor that hash to joining node
}

func (r *Router) onPredecessorLeave() {
	r.lock.Lock()
	defer r.lock.Unlock()
	// move regions that hashed to node that left to successor
}