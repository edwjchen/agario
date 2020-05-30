package region 

import (
	"golang.org/x/net/context"
	// "github.com/paulmach/orb/quadtree"
	// "peer_to_peer/server/player"
	"sync"
	// . "peer_to_peer/server/player_pb"
	"peer_to_peer/server/region_pb"
)

type RegionHandler struct{
	Regions map[uint32]*RegionInfo
	mux     sync.RWMutex
}

func (rh *RegionHandler) Ping(ctx context.Context, request *region_pb.EmptyRequest) (*region_pb.EmptyResponse, error) {
	response := region_pb.EmptyResponse{}
	return &response, nil
}

func (rh *RegionHandler) AddRegion(ctx context.Context, request *region_pb.AddRegionRequest) (*region_pb.EmptyResponse, error) {
	response := region_pb.EmptyResponse{}
	return &response, nil
}

func (rh *RegionHandler) GetRegion(ctx context.Context, request *region_pb.IdRegionRequest) (*region_pb.GetRegionResponse, error) {
	rh.mux.RLock()
	regionId := request.GetId()
	region, _ := rh.Regions[regionId]
	rh.mux.RUnlock()

	allPlayers := make(map[string]*region_pb.Blob)
	for name, p := range region.GetSeen() {
		allPlayers[name] = p.GetBlob()
	}
	for name, p := range region.GetIn() {
		allPlayers[name] = p.GetBlob()
	}
	blobs := []*region_pb.Blob{}
	for _, blob := range allPlayers {
		blobs = append(blobs, blob)
	}

	response := region_pb.GetRegionResponse{
		Blobs:     blobs,
		Foods:     region.GetFood(),
	}
	return &response, nil
}

func (rh *RegionHandler) RemoveRegion(ctx context.Context, request *region_pb.IdRegionRequest) (*region_pb.EmptyResponse, error) {
	//send quit channel
	//close quit channel 
	
	response := region_pb.EmptyResponse{}
	return &response, nil
}

func (rh *RegionHandler) ClientUpdate(ctx context.Context, request *region_pb.UpdateRegionRequest) (*region_pb.UpdateRegionResponse, error) {
	regionId := request.GetId()
	rh.mux.RLock()
	region, _ := rh.Regions[regionId]
	rh.mux.RUnlock()
	
	updatedBlob := request.GetBlob()
	if !updatedBlob.Alive {
		// Remove blob from cache
		region.mux.Lock()
		delete(region.PlayersIn, updatedBlob.Name)
		delete(region.PlayersSeen, updatedBlob.Name)
		region.mux.Unlock()
		response := region_pb.UpdateRegionResponse {
			DeltaMass: 0,
			Alive: false,
		}
		return &response, nil
	}

	// 	If player within region:
	// 		isAlive, eater = didPlayerGetEatenDied()
	// 		If isAlive:
	// 			massIncrease = getMassIncreaseFromFood(player)
	// 			UpdatePos in local region hint cache
	// 		Else
	// 			Remove from cache and do nothing
	// 			eaterServer = routingService.Get(eater)   // eater server is the player server
	// 			eaterServer.MassIncrement(massIncrease)  
	// 	Else:
	// 		massIncrease = getMassIncreaseFromFood(player)
	// 		UpdatePos in local region hint cache

	// if region.BlobIsIn(updatedBlob) {

	// }


	response := region_pb.UpdateRegionResponse{}
	return &response, nil
}

func (rh *RegionHandler) AddFoods(ctx context.Context, request *region_pb.FoodRequest) (*region_pb.EmptyResponse, error) {
	response := region_pb.EmptyResponse{}
	return &response, nil
}

func (rh *RegionHandler) RemoveFoods(ctx context.Context, request *region_pb.FoodRequest) (*region_pb.EmptyResponse, error) {
	response := region_pb.EmptyResponse{}
	return &response, nil
}
