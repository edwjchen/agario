package region

import (
	"golang.org/x/net/context"
	// "github.com/paulmach/orb/quadtree"
	"encoding/binary"
	"hash/fnv"
	"log"
	. "peer_to_peer/common"
	. "peer_to_peer/server/player"
	. "peer_to_peer/server/player_pb"
	. "peer_to_peer/server/region_pb"
	"peer_to_peer/server/router"
	"sync"
	"time"
)

type RegionHandler struct {
	RegionHash    map[uint32]uint32
	Regions       map[uint32]*RegionInfo
	BackupRegions map[uint32]*RegionInfo
	Router        *router.Router
	mux           sync.RWMutex
	RegionChange  chan router.RegionChangeInfo 
}

func getHash(id uint32) uint32 {
	hasher := fnv.New32a()
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, id)
	hasher.Write(b)
	h := hasher.Sum32()
	return h
}

func (rh *RegionHandler) Init() {

	rh.Regions = make(map[uint32]*RegionInfo)
	rh.BackupRegions = make(map[uint32]*RegionInfo)
	rh.RegionHash = make(map[uint32]uint32)

	go rh.Router.Heartbeat()
	time.Sleep(time.Second * 3)

	rh.mux.Lock()

	var i, j uint32
	for i = 0; i < Conf.NREGION_WIDTH; i++ {
		for j = 0; j < Conf.NREGION_HEIGHT; j++ {
			
			regionID := getRegionID(uint16(i), uint16(j))
			h := getHash(regionID)

			rh.RegionHash[regionID] = h
			regionSuccessorHash := rh.Router.Successor(h)

			if regionSuccessorHash == rh.Router.Hash {
				// log.Println("MY REGION!")
				newRegion := &RegionInfo{}
				newRegion.InitRegion(i, j, rh.Router)
				newRegion.SetReady()
				go newRegion.MaintainRegion()
				rh.Regions[regionID] = newRegion
			} else if rh.Router.Successor(regionSuccessorHash+1) == rh.Router.Hash {
				// log.Println("I'M BACKUP!")
				newRegion := &RegionInfo{}
				newRegion.InitRegion(i, j, rh.Router)

				// go newRegion.MaintainRegion()
				rh.BackupRegions[regionID] = newRegion
			} else {
				// log.Println("NOT MY THING!")
			}
		}
	}
	rh.Router.UpdateRingPos()
	rh.mux.Unlock()
	go rh.NodeChangeHandler()
}

func (rh *RegionHandler) Join() {
	// rh.Router.Heartbeat()
	
	rh.Regions = make(map[uint32]*RegionInfo)
	rh.BackupRegions = make(map[uint32]*RegionInfo)
	rh.RegionHash = make(map[uint32]uint32)

	rh.mux.Lock()
	go rh.Router.Heartbeat()
	time.Sleep(time.Second * 1)

	var i, j uint32
	for i = 0; i < Conf.NREGION_WIDTH; i++ {
		for j = 0; j < Conf.NREGION_HEIGHT; j++ {
			
			regionID := getRegionID(uint16(i), uint16(j))
			h := getHash(regionID)

			rh.RegionHash[regionID] = h
		}
	}
	rh.Router.UpdateRingPos()
	rh.mux.Unlock()
	// log.Println("Join unlock")
	go rh.NodeChangeHandler()
	
}

func (rh *RegionHandler) Ping(ctx context.Context, request *EmptyRequest) (*EmptyResponse, error) {
	response := EmptyResponse{}
	return &response, nil
}

func (rh *RegionHandler) GetRegion(ctx context.Context, request *IdRegionRequest) (*GetRegionResponse, error) {
	regionId := request.GetId()
	rh.mux.RLock()
	region, ok := rh.Regions[regionId]
	rh.mux.RUnlock()
	if !ok {
		rh.mux.RLock()
		region, ok = rh.BackupRegions[regionId]
		rh.mux.RUnlock()
		if !ok {
			// log.Println(GetRegionX(regionId), GetRegionY(regionId), "NOT IN PRIMARY NOR BACKUP!")
			return &GetRegionResponse{Blobs: []*Blob{}, Foods: []*Food{}, Ready: false}, nil
		}
	}

	regionReady := region.GetReady()

	// if !regionReady {
	// 	return &GetRegionResponse{Blobs: []*Blob{}, Foods: []*Food{}}, nil
	// }

	allPlayers := make(map[string]*Blob)
	for name, p := range region.GetSeen() {
		allPlayers[name] = p.GetBlob()
	}

	blobs := []*Blob{}
	for _, blob := range allPlayers {
		// log.Println("seq: ", blob.Seq)
		if blob.Alive && blob.Seq > 0 {
			blobs = append(blobs, blob)
		}
	}

	response := GetRegionResponse{
		Blobs: blobs,
		Foods: region.GetFood(),
		Ready: regionReady,
	}
	return &response, nil
}

func (rh *RegionHandler) AddRegions(ctx context.Context, request *AddRegionsRequest) (*EmptyResponse, error) {

	rh.doAddRegions(request.GetRegions())
	// regionId := request.GetId()

	// rh.mux.RLock()
	// log.Println(regionId, "(",rh.Router.Successor(rh.RegionHash[regionId]),",",rh.Router.Successor(rh.Router.Successor(rh.RegionHash[regionId])+1),") addRegionRPC, this is", rh.Router.Hash)
	// region, ok := rh.BackupRegions[regionId]
	// rh.mux.RUnlock()
	// if !ok {
	// 	//log.Println("AddFood race...")
	// 	rh.mux.Lock()
	// 	newRegion := &RegionInfo{}
	// 	newRegion.InitRegion(uint32(GetRegionX(regionId)), uint32(GetRegionY(regionId)), rh.Router)
	// 	newRegion.AddFoods(request.GetFoods())
	// 	rh.BackupRegions[regionId] = newRegion
	// 	rh.mux.Unlock()
	// } else {
	// 	region.AddFoods(request.GetFoods())
	// }

	response := EmptyResponse{}
	return &response, nil
}

func (rh *RegionHandler) doAddRegions(regions []*FoodRequest) {

	for _, request := range regions {
		regionId := request.GetId()

		rh.mux.RLock()
		// log.Println(regionId, "(",rh.Router.Successor(rh.RegionHash[regionId]),",",rh.Router.Successor(rh.Router.Successor(rh.RegionHash[regionId])+1),") addRegionRPC, this is", rh.Router.Hash)
		region, ok := rh.BackupRegions[regionId]
		rh.mux.RUnlock()
		if !ok {
			//log.Println("AddFood race...")
			rh.mux.Lock()
			newRegion := &RegionInfo{}
			newRegion.InitRegion(uint32(GetRegionX(regionId)), uint32(GetRegionY(regionId)), rh.Router)
			newRegion.AddFoods(request.GetFoods())
			rh.BackupRegions[regionId] = newRegion
			rh.mux.Unlock()
		} else {
			region.AddFoods(request.GetFoods())
		}
	}

}

func (rh *RegionHandler) RemoveRegions(ctx context.Context, request *RemoveRegionsRequest) (*EmptyResponse, error) {

	rh.mux.Lock()
	for _, req := range request.GetRegions() {
		delete(rh.BackupRegions, req.GetId())
	}
	rh.mux.Unlock()

	response := EmptyResponse{}
	return &response, nil
}

func (rh *RegionHandler) TransferPrimary(ctx context.Context, request *AddRegionsRequest) (*EmptyResponse, error) {

	rh.doTransferPrimary(request.GetRegions())
	// regionId := request.GetId()

	// rh.mux.Lock()
	// log.Println(regionId, "(",rh.Router.Successor(rh.RegionHash[regionId]),",",rh.Router.Successor(rh.Router.Successor(rh.RegionHash[regionId])+1),") TransferPrimaryRPC, this is", rh.Router.Hash)
	// newRegion := &RegionInfo{}
	// newRegion.InitRegion(uint32(GetRegionX(regionId)), uint32(GetRegionY(regionId)), rh.Router)
	// newRegion.AddFoods(request.GetFoods())
	// rh.Regions[regionId] = newRegion
	// newRegion.SetReady()
	// go newRegion.MaintainRegion()
	// rh.mux.Unlock()

	return &EmptyResponse{}, nil

}

func (rh *RegionHandler) doTransferPrimary(regions []*FoodRequest) {
	for _, request := range regions {
		regionId := request.GetId()
		rh.mux.Lock()
		// log.Println(regionId, "(",rh.Router.Successor(rh.RegionHash[regionId]),",",rh.Router.Successor(rh.Router.Successor(rh.RegionHash[regionId])+1),") TransferPrimaryRPC, this is", rh.Router.Hash)
		newRegion := &RegionInfo{}
		newRegion.InitRegion(uint32(GetRegionX(regionId)), uint32(GetRegionY(regionId)), rh.Router)
		newRegion.AddFoods(request.GetFoods())
		rh.Regions[regionId] = newRegion
		newRegion.SetReady()
		go newRegion.MaintainRegion()
		rh.mux.Unlock()
	}
}

func (rh *RegionHandler) ClientUpdate(ctx context.Context, request *UpdateRegionRequest) (*UpdateRegionResponse, error) {
	// log.Println("Calling ClientUpdate")
	regionId := request.GetId()
	rh.mux.RLock()
	region, ok := rh.Regions[regionId]
	rh.mux.RUnlock()
	if !ok {
		rh.mux.RLock()
		region, ok = rh.BackupRegions[regionId]
		rh.mux.RUnlock()
		if !ok {
			// log.Println(GetRegionX(regionId), GetRegionY(regionId), "NOT IN PRIMARY NOR BACKUP!")
			return &UpdateRegionResponse{DeltaMass: 0, Alive: false, Ready: false}, nil
		}
	}

	updatedBlob := request.GetBlob()

	updatedPlayerInfo := &PlayerInfo{
		Blob:       *updatedBlob,
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
	regionReady := region.GetReady()

	// could use finer granularity locking?...
	region.PlayerSeenMux.Lock()
	defer region.PlayerSeenMux.Unlock()

	existingPlayer, ok := region.PlayersSeen[updatedBlobID]
	if ok {
		// blob already dead: DEAD
		if !region.PlayersSeen[updatedBlobID].Blob.Alive {
			return &UpdateRegionResponse{DeltaMass: 0, Alive: false, Ready: regionReady}, nil
		} 
		// OO info: ignore
		if updatedBlob.Seq > 0 && existingPlayer.Blob.Seq > updatedBlob.Seq {
			return &UpdateRegionResponse{DeltaMass: 0, Alive: true, Ready: regionReady}, nil
		}
	} 

	region.PlayersSeen[updatedBlobID] = updatedPlayerInfo

	if !updatedBlob.Alive {
		// Remove blob from cache
		return &UpdateRegionResponse{DeltaMass: 0, Alive: false, Ready: regionReady}, nil
	}

	// Eviction info: just ignore
	if updatedBlob.Seq < 0 {
		return &UpdateRegionResponse{DeltaMass: 0, Alive: true, Ready: regionReady}, nil
	}

	//if !regionReady {
	//	return &UpdateRegionResponse{DeltaMass: 0, Alive: true, Ready: regionReady}, nil
	//}

	var massIncrease int32

	if region.BlobIsIn(updatedBlob) {
		isAlive, eater := region.WasEaten(updatedBlob)
		if isAlive {
			massIncrease = region.GetNumFoodsEaten(updatedBlob)
		} else {
			updatedBlob.Alive = false
			// region.PlayersIn[updatedBlobID].Blob.Alive = false
			region.PlayersSeen[updatedBlobID].Blob.Alive = false

			eaterServer := rh.Router.GetPlayerConn(eater.Ip)
			if eaterServer != nil {
				client := NewPlayerClient(eaterServer)
				massIncReq := &MassIncrementRequest{MassIncrease: updatedBlob.Mass}
				_, err := client.MassIncrement(context.Background(), massIncReq)
				if err != nil {
					// log.Println("Failed increment", err)
					rh.Router.InvalidatePlayerConn(eater.Ip)
				}
			}
			// log.Println("DEAD!!!")
			return &UpdateRegionResponse{DeltaMass: 0, Alive: false, Ready: regionReady}, nil
		}
	} else {
		massIncrease = region.GetNumFoodsEaten(updatedBlob)
	}

	response := UpdateRegionResponse{
		DeltaMass: massIncrease,
		Alive:     updatedBlob.Alive,
		Ready:     true,
	}
	return &response, nil
}

// below two methods are for replication
func (rh *RegionHandler) AddFoods(ctx context.Context, request *FoodRequest) (*EmptyResponse, error) {
	
	regionId := request.GetId()
	//log.Println(regionId, "(",rh.Router.Successor(rh.RegionHash[regionId]),",",rh.Router.Successor(rh.Router.Successor(rh.RegionHash[regionId])+1),") addFoodsRPC, this is", rh.Router.Hash)
	rh.mux.RLock()
	region, ok := rh.BackupRegions[regionId]
	rh.mux.RUnlock()
	if !ok {
		//log.Println("AddFood race...")
		rh.mux.Lock()
		newRegion := &RegionInfo{}
		newRegion.InitRegion(uint32(GetRegionX(regionId)), uint32(GetRegionY(regionId)), rh.Router)
		newRegion.AddFoods(request.GetFoods())
		rh.BackupRegions[regionId] = newRegion
		rh.mux.Unlock()
	} else {
		region.AddFoods(request.GetFoods())
	}

	response := EmptyResponse{}
	return &response, nil
}

func (rh *RegionHandler) RemoveFoods(ctx context.Context, request *FoodRequest) (*EmptyResponse, error) {
	
	regionId := request.GetId()
	//log.Println(regionId, "(",rh.Router.Successor(rh.RegionHash[regionId]),",",rh.Router.Successor(rh.Router.Successor(rh.RegionHash[regionId])+1),") rmFoodsRPC, this is", rh.Router.Hash)
	rh.mux.RLock()
	region, ok := rh.BackupRegions[regionId]
	rh.mux.RUnlock()
	if ok {
		region.RemoveFoods(request.GetFoods())
	}

	response := EmptyResponse{}
	return &response, nil
}

func GetRegionX(rid uint32) uint16 {
	return uint16((rid & (0xffff0000)) >> 16)
}
func GetRegionY(rid uint32) uint16 {
	return uint16((rid & (0x0000ffff)))
}

// rule of thumb: never do an RPC while holding a lock!
func (rh *RegionHandler) NodeChangeHandler() {
	// log.Println("NodeChangeHandler started")
	for {
		ncInfo := <-rh.RegionChange

		if ncInfo.PrevConn == nil && ncInfo.CurrConn == nil {

			// one node edge case: dump everything into primary
			// log.Println("one node edge case")
			rh.mux.Lock()
			for rid, r := range rh.BackupRegions {
				rh.Regions[rid] = r 
				r.SetReady()
				go r.MaintainRegion()
				delete(rh.BackupRegions, rid)
			}
			// log.Println("Number of reg handling: ")
			rh.mux.Unlock()


		} else if ncInfo.Successor && ncInfo.Join {

			// onSuccessorJoin:
			// move regions on curr node that hash to curr region to joined node
			// log.Println("onSuccessorJoin", ncInfo.Curr, ncInfo.Prev)
			regionsCopy := []*FoodRequest{}
			rh.mux.RLock()
			for rid, r := range rh.Regions {
				regionsCopy = append(regionsCopy, &FoodRequest{Id:rid, Foods:r.GetFood()}) 
			}
			rh.mux.RUnlock()

			conn := rh.Router.GetSuccessor()
			regionClient := NewRegionClient(conn)
			_, err := regionClient.AddRegions(context.Background(), &AddRegionsRequest{Regions: regionsCopy})
			if err != nil {
				log.Println("Successor Join big no no: ", err)
			}
			// for rid, r := range regionsCopy {
			// 	regionClient := NewRegionClient(conn)
			// 	_, err := regionClient.AddRegions(context.Background(), &AddRegionRequest{Id: rid, Foods: r.GetFood()})
			// 	if err != nil {
			// 		log.Println("Successor Join big no no: ", err)
			// 	}
			// 	// conn := ncInfo.PrevConn
			// 	// _, err = regionClient.RemoveRegion(context.Background(), &region_pb.IdRegionRequest{Id: rid})
			// 	// if err != nil {
			// 	// 	log.Println("Successor Join big no no: ", err)
			// 	// }
			// }

		} else if ncInfo.Successor && !ncInfo.Join {

			// onSuccessorLeave:
			// move regions on curr node that hash to curr region to new successor
			// log.Println("onSuccessorLeave", ncInfo.Curr, ncInfo.Prev)

			regionsCopy := []*FoodRequest{}
			rh.mux.RLock()
			for rid, r := range rh.Regions {
				regionsCopy = append(regionsCopy, &FoodRequest{Id:rid, Foods:r.GetFood()}) 
			}
			rh.mux.RUnlock()

			// regionsCopy := make(map[uint32]*RegionInfo)
			// rh.mux.RLock()
			// for rid, r := range rh.Regions {
			// 	regionsCopy[rid] = r 
			// }
			// rh.mux.RUnlock()

			conn := rh.Router.GetSuccessor()
			regionClient := NewRegionClient(conn)
			_, err := regionClient.AddRegions(context.Background(), &AddRegionsRequest{Regions: regionsCopy})
			if err != nil {
				log.Println("Successor Join big no no: ", err)
			}
			// for rid, r := range regionsCopy {
			// 	regionClient := NewRegionClient(conn)
			// 	_, err := regionClient.AddRegion(context.Background(), &AddRegionRequest{Id: rid, Foods: r.GetFood()})
			// 	if err != nil {
			// 		log.Println("Successor Join big no no: ", err)
			// 	}
			// }

		} else if !ncInfo.Successor && ncInfo.Join {

			// on predecessor join:
			// remove regions on curr node that hash to prepredecessor
			// log.Println("on predecessor join", ncInfo.Curr, ncInfo.Prev)
			rh.mux.Lock()
			for rid, _ := range rh.BackupRegions {
				if rh.Router.Successor(rh.RegionHash[rid]) != ncInfo.Curr {
					delete(rh.BackupRegions, rid)
				}
			}
			rh.mux.Unlock()

			// remove region on successor for which now I'm backup
			// regionsCopy := make(map[uint32]*RegionInfo)
			regionsCopy := []*FoodRequest{}
			rmRegions  := []*IdRegionRequest{}
			rh.mux.Lock()
			for rid, r := range rh.Regions {
				log.Println("Checking", rid, rh.RegionHash[rid], rh.Router.Successor(rh.RegionHash[rid]), ncInfo.Curr, rh.Router.Hash)
				if rh.Router.Successor(rh.RegionHash[rid]) == ncInfo.Curr {
					// regionsCopy[rid] = r
					regionsCopy = append(regionsCopy, &FoodRequest{Id: rid, Foods:r.GetFood()}) 
					rmRegions   = append(rmRegions, &IdRegionRequest{Id: rid})
					r.Quit <- true
					r.ClearAllBlobCache()
					rh.BackupRegions[rid] = r
					delete(rh.Regions, rid)
					// r.UnsetReady()
				}
			}
			rh.mux.Unlock()

			// move regions that hash to joining node from curr node to joining node
			// remove regions on successor that hash to joining node
			successorConn  := rh.Router.GetSuccessor()
			PredcessorConn := rh.Router.GetPredecessor()
			regionClient := NewRegionClient(successorConn)
			_, err := regionClient.RemoveRegions(context.Background(), &RemoveRegionsRequest{Regions: rmRegions})
			if err != nil {
				log.Println("Predecessor Join big no no: ", err)
			}
			regionClient = NewRegionClient(PredcessorConn)
			_, err = regionClient.TransferPrimary(context.Background(), &AddRegionsRequest{Regions: regionsCopy})
			if err != nil {
				log.Println("Predecessor Join big no no: ", err)
			}

			// for rid, r := range regionsCopy {
			// 	regionClient := NewRegionClient(successorConn)
			// 	_, err := regionClient.RemoveRegion(context.Background(), &IdRegionRequest{Id: rid})
			// 	if err != nil {
			// 		log.Println("Predecessor Join big no no: ", err)
			// 	}

			// 	regionClient = NewRegionClient(PredcessorConn)
			// 	_, err = regionClient.TransferPrimary(context.Background(), &FoodRequest{Id: rid, Foods: r.GetFood()})
			// 	if err != nil {
			// 		log.Println("Predecessor Join big no no: ", err)
			// 	}
			// 	r.UnsetReady()
			// }

		} else {

			// move regions that hashed to node that left to successor
			// log.Println("on predecessor leave", ncInfo.Curr, ncInfo.Prev)

			regionsCopy := []*FoodRequest{}
			// regionsCopy := make(map[uint32]*RegionInfo)

			rh.mux.Lock()
			for rid, r := range rh.BackupRegions {
				// log.Println("Checking", rid, rh.RegionHash[rid])
				if rh.Router.Successor(rh.RegionHash[rid]) == rh.Router.Hash {
					// log.Println("Moving", rid, "from bkup to prim")
					regionsCopy = append(regionsCopy, &FoodRequest{Id: rid, Foods:r.GetFood()})
					//regionsCopy[rid] = r
					rh.Regions[rid] = r
					r.SetReady()
					go r.MaintainRegion()
					delete(rh.BackupRegions, rid)
				}
			}
			rh.mux.Unlock()

			successorConn := rh.Router.GetSuccessor()
			regionClient := NewRegionClient(successorConn)
			_, err := regionClient.AddRegions(context.Background(), &AddRegionsRequest{Regions: regionsCopy})
			if err != nil {
				log.Println("Predecessor leave big no no: ", err)
			}
			// for rid, r := range regionsCopy {
			// 	regionClient := NewRegionClient(successorConn)
			// 	_, err := regionClient.AddRegions(context.Background(), &AddRegionRequest{Id: rid, Foods: r.GetFood()})
			// 	if err != nil {
			// 		log.Println("Predecessor leave big no no: ", err)
			// 	}
			// }
		}
		
	}
}

func getRegionID(x, y uint16) uint32 {
	return uint32(x) << 16 | uint32(y)
}
