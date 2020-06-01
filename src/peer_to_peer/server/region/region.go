package region 

import (
	"golang.org/x/net/context"
	// "github.com/paulmach/orb/quadtree"
	"sync"
	"strings"
	"fmt"
	"time"
	"log"
	"peer_to_peer/server/router"
	. "peer_to_peer/server/player"
	. "peer_to_peer/server/player_pb"
	. "peer_to_peer/server/region_pb"
	. "peer_to_peer/common"
)

type RegionHandler struct{
	Regions map[uint32]*RegionInfo
	Router *router.Router
	mux     sync.RWMutex
}

func (rh *RegionHandler) Ping(ctx context.Context, request *EmptyRequest) (*EmptyResponse, error) {
	response := EmptyResponse{}
	return &response, nil
}

func (rh *RegionHandler) GetRegion(ctx context.Context, request *IdRegionRequest) (*GetRegionResponse, error) {
	rh.mux.RLock()
	regionId := request.GetId()
	region, _ := rh.Regions[regionId]
	rh.mux.RUnlock()

	allPlayers := make(map[string]*Blob)
	for name, p := range region.GetSeen() {
		allPlayers[name] = p.GetBlob()
	}
	for name, p := range region.GetIn() {
		allPlayers[name] = p.GetBlob()
	}
	blobs := []*Blob{}
	for _, blob := range allPlayers {
		blobs = append(blobs, blob)
	}

	response := GetRegionResponse{
		Blobs:     blobs,
		Foods:     region.GetFood(),
	}
	return &response, nil
}

func (rh *RegionHandler) AddRegion(ctx context.Context, request *AddRegionRequest) (*EmptyResponse, error) {
	response := EmptyResponse{}
	return &response, nil
}

func (rh *RegionHandler) RemoveRegion(ctx context.Context, request *IdRegionRequest) (*EmptyResponse, error) {
	//send quit channel
	//close quit channel 
	
	response := EmptyResponse{}
	return &response, nil
}

func BlobID(blob *Blob) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s:%x", blob.Ip, blob.Ver)
	return b.String()
}

func (rh *RegionHandler) ClientUpdate(ctx context.Context, request *UpdateRegionRequest) (*UpdateRegionResponse, error) {
	regionId := request.GetId()
	rh.mux.RLock()
	region, _ := rh.Regions[regionId]
	rh.mux.RUnlock()
	
	updatedBlob := request.GetBlob()
	updatedPlayerInfo := &PlayerInfo{
		Blob: *updatedBlob,
		LastUpdate: time.Now(),
	}
	updatedBlobID := BlobID(updatedBlob)

	// bg thread doing time-based cache eviction

	// instead of delete, mark player as dead
	// if player is dead and we receive update saying player is alive, tell sender that player is dead.

	/*
	t=0: player sends update1 to server
	t=1: server received update1 and calculates that player is dead. eagerly evict player from cache
	t=2: player sends update2 to server (alive)
	t=3: server sends response1 to tell player that he's dead
	t=4: server receives update2 (alive). put player into cache
	*/

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
	// 		massIncrease = getMasncreaseFromFood(player)
	// 		UpdatePos in local region hint cache

	region.PlayerInMux.Lock()
	defer region.PlayerInMux.Unlock()
	region.PlayerSeenMux.Lock()
	defer region.PlayerSeenMux.Unlock()

	region.PlayersIn[updatedBlobID] = updatedPlayerInfo
	region.PlayersSeen[updatedBlobID] = updatedPlayerInfo
	
	if !updatedBlob.Alive {
		// Remove blob from cache
		response := UpdateRegionResponse {
			DeltaMass: 0,
			Alive: false,
		}
		return &response, nil
	}

	var massIncrease int32

	if region.BlobIsIn(updatedBlob) {
		isAlive, eater := region.WasEaten(updatedBlob)
		if isAlive {
			massIncrease = region.GetNumFoodsEaten(updatedBlob)
		} else {
			updatedBlob.Alive = false
			region.PlayersIn[updatedBlobID].Blob.Alive = false
			region.PlayersSeen[updatedBlobID].Blob.Alive = false
			
			eaterServer := rh.Router.GetPlayerConn(eater.Ip)
			if eaterServer != nil {
				client := NewPlayerClient(eaterServer)
				massIncReq := &MassIncrementRequest{MassIncrease: updatedBlob.Mass}
				_, err := client.MassIncrement(context.Background(), massIncReq)
				if err != nil {
					log.Println("Failed increment", err)
					rh.Router.InvalidatePlayerConn(eater.Ip)
				}
			}
		}
	} else {
		massIncrease = region.GetNumFoodsEaten(updatedBlob)
	}

	response := UpdateRegionResponse{
		DeltaMass: massIncrease,
		Alive: updatedBlob.Alive,
	}
	return &response, nil
}

// below two methods are for replication
func (rh *RegionHandler) AddFoods(ctx context.Context, request *FoodRequest) (*EmptyResponse, error) {
	response := EmptyResponse{}
	return &response, nil
}

func (rh *RegionHandler) RemoveFoods(ctx context.Context, request *FoodRequest) (*EmptyResponse, error) {
	response := EmptyResponse{}
	return &response, nil
}
