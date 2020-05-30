package router

//client in python <---> player server <---> region server

import (
	"hash/fnv"
	"sync"
	"log"
	"time"
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"peer_to_peer/server/region_pb"
	"sort"
)

type Router struct {
	lock sync.Mutex 
	conns map[string] *grpc.ClientConn// map from ip of ip to connections
	iphash map[string]uint32
	haship map[uint32]string
	liveBacks []uint32 // stores hashes of ip
	hash uint32
}

type heartbeatOutput struct {
	ip   string
	conn *grpc.ClientConn
}

// Need to store hash of self
func (r *Router) Init(servers []string) {
	r.haship = make(map[uint32]string)
	r.iphash = make(map[string]uint32)
	r.conns = make(map[string] *grpc.ClientConn)
	for _, ip := range(servers) {
		hasher := fnv.New32a()
		hasher.Write([]byte(ip))
		hash := uint32(hasher.Sum32())
		r.iphash[ip] = hash 
		r.haship[hash] = ip
		r.conns[ip] = nil
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

		log.Println("OK")
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
func (r *Router) Get(key string) *grpc.ClientConn {
	// return grpc connection of head of chain
	hasher := fnv.New32a()
	hasher.Write([]byte(key))
	hash := uint32(hasher.Sum32())

	primaryHash := r.successor(hash)
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.conns[r.haship[primaryHash]]
}

// Returns GRPC connection
func (r *Router) GetSuccessor() *grpc.ClientConn {
	// get whatever is after us in aliveBacks
	successorHash := r.successor(r.hash)
	r.lock.Lock()
	defer r.lock.Unlock()
	return r.conns[r.haship[successorHash]]
}

func (r *Router) successor(h uint32) uint32 {
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