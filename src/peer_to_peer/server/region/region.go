package region 

import (
	"golang.org/x/net/context"
	"github.com/paulmach/orb/quadtree"
	"peer_to_peer/server/player"
	"sync"
	// "log"
)

type RegionHandler struct{
	regions         map[uint32]*RegionInfo
}

func (RegionHandler) Ping(ctx context.Context, request *EmptyRequest) (*EmptyResponse, error) {
	response := EmptyResponse{}
	return &response, nil
}

