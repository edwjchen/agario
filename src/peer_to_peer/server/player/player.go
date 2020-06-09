package player

import (
	"fmt"
	"log"
	"math"
	"sync"
	. "peer_to_peer/common"
	. "peer_to_peer/server/player_pb"
	. "peer_to_peer/server/region_pb"
	"peer_to_peer/server/router"
	"strings"

	"golang.org/x/net/context"
)

type PlayerHandler struct {
	Player      PlayerInfo
	Mux         sync.Mutex
	Router      *router.Router
	PrevRegions map[uint32]bool
}

func (ph *PlayerHandler) Init(ctx context.Context, request *InitRequest) (*InitResponse, error) {
	ph.Mux.Lock()
	newBlobIp, startX, startY, mass, ver := ph.Player.NewBlob()
	ph.Mux.Unlock()
	response := InitResponse{
		Ip:   newBlobIp,
		Ver:  ver,
		X:    startX,
		Y:    startY,
		Mass: mass,
	}

	return &response, nil
}

func doClientUpdate(regionId uint32, c chan *UpdateRegionResponse, blob *Blob, r *router.Router) {
	// use router to get the grpc clientconn,
	primary, backup := r.Get(regionId)
	// create client stub from clientconn
	regionClient := NewRegionClient(primary)
	clientUpdate := UpdateRegionRequest{Blob: blob, Id: regionId}
	response, err := regionClient.ClientUpdate(context.Background(), &clientUpdate)
	// log.Println(err)
	if err != nil {
		if primary != backup {
			regionClient = NewRegionClient(backup)
			response, err = regionClient.ClientUpdate(context.Background(), &clientUpdate)
			if err != nil {
				log.Fatalln(regionId, "Failed to get response from both prim and bkup!")
			} 
		} else {
			log.Fatalln(regionId, "only has one associated node and is down!")
		}
	} 
	// call method on goroutine
	c <- response
}

// Search function responsible to get the Country information
func (ph *PlayerHandler) Move(ctx context.Context, request *MoveRequest) (*MoveResponse, error) {
	// for now just echo response with increment on position
	// log.Println("Moving!")
	ph.Mux.Lock()
	defer ph.Mux.Unlock()

	if !ph.Player.GetAlive() {
		return &MoveResponse{
			X:     0,
			Y:     0,
			Alive: false,
			Mass:  0,
		}, nil

		 // &response, nil
	}

	dx := request.GetX()
	dy := request.GetY()

	ph.Player.UpdateSeq()

	// log.Println("dw & dy", dx, dy)

	// log.Println("get: ", dx, dy)
	rotation := math.Atan2(dy-float64(Conf.SCREEN_HEIGHT/2), dx-float64(Conf.SCREEN_WIDTH/2)) * 180 / math.Pi
	vx := Conf.SPEED * (90 - math.Abs(rotation)) / 90
	var vy float64
	if rotation < 0 {
		vy = -1*Conf.SPEED + math.Abs(vx)
	} else {
		vy = Conf.SPEED - math.Abs(vx)
	}
	// log.Println("vx & vy", vx, vy)

	x, y := ph.Player.UpdatePos(vx, vy)

	regions := ph.Player.GetAOI() // returns list of region_id

	// log.Println("AOI", regions)
	resChan := make(chan *UpdateRegionResponse, len(regions))
	blob := ph.Player.GetBlob()

	newRegions := make(map[uint32]bool)

	for _, regionId := range regions {
		// log.Println("Adding ", regionId)
		go doClientUpdate(regionId, resChan, blob, ph.Router)
		delete(ph.PrevRegions, regionId)
		newRegions[regionId] = true
	}

	evictBlob := ph.Player.GetBlob()
	evictBlob.Seq = -1
	evictChan := make(chan *UpdateRegionResponse, len(ph.PrevRegions))
	for erid, _ := range ph.PrevRegions {
		// log.Println("Evicting ", erid)
		go doClientUpdate(erid, evictChan, evictBlob, ph.Router)
	}
	ph.PrevRegions = newRegions

	for _ = range regions {
		// log.Println("I got a response!")
		resp := <-resChan
		if !resp.Alive {
			ph.Player.Die()

			evictBlob := ph.Player.GetBlob()
			evictChan := make(chan *UpdateRegionResponse, len(ph.PrevRegions))
			for erid, _ := range ph.PrevRegions {
				// log.Println("removing from ", erid)
				go doClientUpdate(erid, evictChan, evictBlob, ph.Router)
			}
			ph.PrevRegions = make(map[uint32]bool)
			
			return &MoveResponse{
				X:     0,
				Y:     0,
				Alive: false,
				Mass:  0,
			}, nil
		}
		ph.Player.IncrementMass(resp.DeltaMass)
	}

	response := MoveResponse{
		X:     x,
		Y:     y,
		Alive: true,
		Mass:  ph.Player.GetMass(), //TODO: change but leave for now
	}

	return &response, nil
}

func doRegionUpdate(regionId uint32, c chan *GetRegionResponse, r *router.Router) {
	// use router to get the grpc clientconn,
	primary, backup := r.Get(regionId)
	// create client stub from clientconn
	regionClient := NewRegionClient(primary)
	getRegionRequest := IdRegionRequest{Id: regionId}
	response, err := regionClient.GetRegion(context.Background(), &getRegionRequest)
	
	if err != nil {
		if primary != backup {
			regionClient = NewRegionClient(backup)
			response, err = regionClient.GetRegion(context.Background(), &getRegionRequest)
			if err != nil {
				log.Fatalln(regionId, "Failed to get response from both prim and bkup!")
			} 
		} else {
			log.Fatalln(regionId, "only has one associated node and is down!")
		}
	} 

	// call method on goroutine
	c <- response
}

func (ph *PlayerHandler) Region(ctx context.Context, request *RegionRequest) (*RegionResponse, error) {
	// regionId := request.GetId()
	// conn := ph.router.Get(regionId)
	// regionClient := region_pb.NewRegionClient(conn)

	// regionRequest := &region_pb.IdRegionRequest{Id: regionId}
	// getRegionResponse, _ := regionClient.GetRegion(context.Background(), regionRequest)

	// response := region_pb.RegionResponse{
	// 	Blobs: getRegionResponse.Blobs,
	// 	Foods: getRegionResponse.Foods,
	// }

	//Get Info from Regions
	//Compile Info from Regions
	visibleBlobs := make(map[string]*Blob)
	visibleFoods := make(map[*Food]bool)

	ph.Mux.Lock()
	regions := ph.Player.GetAOI() // returns list of region_id
	ph.Mux.Unlock()

	resChan := make(chan *GetRegionResponse, len(regions))

	for _, regionId := range regions {
		doRegionUpdate(regionId, resChan, ph.Router)
	}

	for _= range regions {
		response := <-resChan
		blobs := response.GetBlobs()
		foods := response.GetFoods()

		for _, b := range blobs {
			bid := BlobID(b)
			// log.Println("Region", reg, "got blob ", b)
			if existingBlob, exists := visibleBlobs[bid]; exists {
				if b.Seq > existingBlob.Seq {
					visibleBlobs[bid] = b
				}
			} else {
				visibleBlobs[bid] = b
			}
		}

		for _, f := range foods {
			visibleFoods[f] = true
		}
	}

	//compile list of blobs and foods
	retBlobs := make([]*Blob, 0)
	retFoods := make([]*Food, 0)

	for _, b := range visibleBlobs {
		retBlobs = append(retBlobs, b)
	}

	for f, _ := range visibleFoods {
		retFoods = append(retFoods, f)
	}

	//Respond]
	response := RegionResponse{
		Blobs: retBlobs,
		Foods: retFoods,
	}

	return &response, nil
}

func (ph *PlayerHandler) MassIncrement(ctx context.Context, request *MassIncrementRequest) (*MassIncrementResponse, error) {
	deltaMass := request.GetMassIncrease()
	ph.Mux.Lock()
	ph.Player.IncrementMass(deltaMass)
	ph.Mux.Unlock()

	response := MassIncrementResponse{}
	return &response, nil
}

func GetRadiusFromMass(mass int32) float64 {
	rad := math.Sqrt(float64(mass)) * float64(Conf.MASS_MULTIPLIER)
	return rad
}

func BlobID(blob *Blob) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%s:%x", blob.Ip, blob.Ver)
	return b.String()
}
