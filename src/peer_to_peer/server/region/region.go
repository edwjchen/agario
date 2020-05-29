package region 

import (
	"golang.org/x/net/context"
	"github.com/paulmach/orb/quadtree"
	"sync"
	// "log"
)

type RegionHandler struct{
	regions         map[uint32]*Region
	playersBelongs  map[string]*Player
	playersPresents map[string]*Player
}

type Region struct {
	foodTree *quadtree.Quadtree
	mux      sync.Mutex
	x        uint16   
	y        uint16   
	hash     uint32   
}

func (RegionHandler) Ping(ctx context.Context, request *EmptyRequest) (*EmptyResponse, error) {
	response := EmptyResponse{}
	return &response, nil
}

