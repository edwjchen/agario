package player

import (
	"golang.org/x/net/context"
	"math"
	"peer_to_peer/server/router"
	. "peer_to_peer/common"
	. "peer_to_peer/server/region_pb"
	. "peer_to_peer/server/player_pb"
	// "log"
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
	rotation := math.Atan2(dy-SCREEN_HEIGHT/2, dx-SCREEN_WIDTH/2) * 180 / math.Pi
	vx := SPEED * (90 - math.Abs(rotation)) / 90
	var vy float64
	if rotation < 0 {
		vy = -1*SPEED+ math.Abs(vx)
	} else {
		vy = SPEED - math.Abs(vx)
	}
	// log.Println("vx & vy", vx, vy)

	x, y := ph.Player.UpdatePos(vx, vy)

	response := MoveResponse{
		X:     x,
		Y:     y,
		Alive: true,
		Mass:  ph.Player.GetMass(),//TODO: change but leave for now
	}

	regions := ph.Player.GetAOI() // returns list of region_id
	resChan := make(chan bool)
	blob := ph.Player.GetBlob()
	regionCall := func(regionId uint32, c chan bool) {
		// use router to get the grpc clientconn, 
		conn := ph.Router.Get(regionId)
		// create client stub from clientconn
		regionClient := NewRegionClient(conn)
		clientUpdate := UpdateRegionRequest{Blob: blob, Id: regionId}
		regionClient.ClientUpdate( context.Background(), &clientUpdate )

		// call method on goroutine
		c <- true
	}

	for _, regionId := range(regions) {
		regionCall(regionId, resChan)
	}

	for _ = range(regions) {
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

	response := RegionResponse{
		Blobs: []*Blob{ph.Player.GetBlob()},
		Foods: make([]*Food, 0),
	}

	return &response, nil
}

func (ph *PlayerHandler) MassIncrement(ctx context.Context, request *MassIncrementRequest) (*MassIncrementResponse, error) {
	deltaMass := request.GetMassIncrease()
	ph.Player.IncrementMass(deltaMass)

	return nil, nil
}