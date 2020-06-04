package player

import (
	"log"
	"math"
	. "peer_to_peer/common"
	. "peer_to_peer/server/player_pb"
	. "peer_to_peer/server/region_pb"
	"peer_to_peer/server/router"

	"golang.org/x/net/context"
)

type PlayerHandler struct {
	Player PlayerInfo
	Router *router.Router
}

func (ph *PlayerHandler) Init(ctx context.Context, request *InitRequest) (*InitResponse, error) {
	newBlobIp, startX, startY, mass, ver := ph.Player.NewBlob()
	response := InitResponse{
		Ip:   newBlobIp,
		Ver:  ver,
		X:    startX,
		Y:    startY,
		Mass: mass,
	}

	return &response, nil
}

// Search function responsible to get the Country information
func (ph *PlayerHandler) Move(ctx context.Context, request *MoveRequest) (*MoveResponse, error) {
	// for now just echo response with increment on position
	// log.Println("Moving!")
	if !ph.Player.GetAlive() {
		response := MoveResponse{
			X:     0,
			Y:     0,
			Alive: false,
			Mass:  0,
		}

		return &response, nil
	}

	dx := request.GetX()
	dy := request.GetY()

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

	response := MoveResponse{
		X:     x,
		Y:     y,
		Alive: true,
		Mass:  ph.Player.GetMass(), //TODO: change but leave for now
	}

	regions := ph.Player.GetAOI() // returns list of region_id
	log.Println("AOI", regions)
	resChan := make(chan bool)
	blob := ph.Player.GetBlob()
	regionCall := func(regionId uint32, c chan bool) {
		// use router to get the grpc clientconn,
		conn := ph.Router.Get(regionId)
		// create client stub from clientconn
		regionClient := NewRegionClient(conn)
		clientUpdate := UpdateRegionRequest{Blob: blob, Id: regionId}
		_, err := regionClient.ClientUpdate(context.Background(), &clientUpdate)
		log.Println(err)
		if err != nil {
			log.Println("client updates big no no: ", err)
		}
		// call method on goroutine
		c <- true
	}

	for _, regionId := range regions {
		go regionCall(regionId, resChan)
	}

	for _ = range regions {
		log.Println("I got a response!")
		<-resChan
	}

	return &response, nil
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
	visibleBlobs := make(map[*Blob]bool)
	visibleFoods := make(map[*Food]bool)

	regions := ph.Player.GetAOI() // returns list of region_id
	resChan := make(chan *GetRegionResponse)
	regionCall := func(regionId uint32, c chan *GetRegionResponse) {
		// use router to get the grpc clientconn,
		conn := ph.Router.Get(regionId)
		// create client stub from clientconn
		regionClient := NewRegionClient(conn)
		getRegionRequest := IdRegionRequest{Id: regionId}
		response, _ := regionClient.GetRegion(context.Background(), &getRegionRequest)

		// call method on goroutine
		c <- response
	}

	for _, regionId := range regions {
		go regionCall(regionId, resChan)
	}

	for _ = range regions {
		response := <-resChan
		blobs := response.GetBlobs()
		foods := response.GetFoods()

		// TODO: Fix dedup
		for _, b := range blobs {
			visibleBlobs[b] = true
		}

		for _, f := range foods {
			visibleFoods[f] = true
		}
	}

	//compile list of blobs and foods
	retBlobs := make([]*Blob, 0)
	retFoods := make([]*Food, 0)

	for b, _ := range visibleBlobs {
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
	ph.Player.IncrementMass(deltaMass)

	response := MassIncrementResponse{}
	return &response, nil
}
