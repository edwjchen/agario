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

var RegionInfoStruct RegionInfo


func (RegionHandler) Ping(ctx context.Context, request *EmptyRequest) (*EmptyResponse, error) {
	response := EmptyResponse{}
	return &response, nil
}


func (RegionHandler) AddRegion(ctx context.Context, request *AddRegionRequest) (*EmptyResponse, error) {
	response := EmptyResponse{}
	return &response, nil
}

func (RegionHandler) GetRegion(ctx context.Context, request *IdRegionRequest) (*GetRegionResponse, error) {
	response := GetRegionResponse{
		Blobs:     make([]*Blob, 0),
		Foods:     make([]*Food, 0),
	}
	return &response, nil
}

func (RegionHandler) RemoveRegion(ctx context.Context, request *IdRegionRequest) (*EmptyResponse, error) {
	response := EmptyResponse{}
	return &response, nil
}

func (RegionHandler) RemoveRegion(ctx context.Context, request *IdRegionRequest) (*EmptyResponse, error) {
	response := EmptyResponse{}
	return &response, nil
}

func (RegionHandler) ClientUpdate(ctx context.Context, request *UpdateRegionRequest) (*UpdateRegionResponse, error) {
	response := UpdateRegionResponse{}
	return &response, nil
}

func (RegionHandler) AddFoods(ctx context.Context, request *FoodRequest) (*EmptyResponse, error) {
	response := EmptyResponse{}
	return &response, nil
}

func (RegionHandler) RemoveFoods(ctx context.Context, request *FoodRequest) (*EmptyResponse, error) {
	response := EmptyResponse{}
	return &response, nil
}
