package server

//client in python <---> player server <---> region server

import (
	"hash/fnv"
	"sync"
	"time"
	"google.golang.org/grpc"
	"peer_to_peer/server/region"
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

func (*r Router) init(servers []string) {
	
}

// this heartbeat function runs as a bg goroutine
func (*r Router) heartbeat {

	retChan := make(chan heartbeatOutput)
	ping := func(ip string, cxn *grpc.ClientConn) {
		if cxn == nil {
			conn, err := grpc.Dial(ip)
		}

		retChan <- {
			ip:   ip,
			conn: cxn,
		}
	}

	for {
		<-time.Tick(time.Second)
		lock.Lock()
		for h, cxn := range r.conns {
			// send rpc
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

		sort.Slice(self.AliveBacks, func(i, j int) bool {
		  	return r.liveBacks[i] < self.liveBacks[j]
		})

		lock.Unlock()
	}
}

// Returns GRPC connection
func (*r Router) Get(key string) {
	// return grpc connection of head of chain
	hasher := fnv.New32a()
	hasher.Write([]byte(key))
	hash := uint32(hasher.Sum32())

	primaryHash := r.successor(hash)
	r.Lock()
	return r.conns[r.haship[primaryHash]}
	defer r.Unlock()
}

// Returns GRPC connection
func (*r Router) GetSuccessor() {
	// get whatever is after us in aliveBacks
	successorHash := r.successor(r.hash)
	r.Lock()
	return r.conns[r.haship[successorHash]]
	defer r.Unlock()
}

func (*r Router) successor(h uint32) uint32 {
	r.Lock()
	defer r.Unlock()
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

func (*r Router) onSuccessorJoin {
	r.Lock()
	defer r.Unlock()
	// move regions on curr node that hash to curr region to joined node
}

func (*r Router) onSuccessorLeave {
	r.Lock()
	defer r.Unlock()
	// move regions on curr node that hash to curr region to new successor
}

func (*r Router) onPredecessorJoin {
	r.Lock()
	defer r.Unlock()
	// remove regions on curr node that hash to prepredecessor
}

func (*r Router) onPredecessorLeave {
	r.Lock()
	defer r.Unlock()
	// move regions that hashed to node that left to successor
}