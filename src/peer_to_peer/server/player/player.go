package player

import (
	"golang.org/x/net/context"
	"math"
	"peer_to_peer/server/router"
	"peer_to_peer/server/player_pb"
	// "peer_to_peer/server/region_pb"
	// "log"
)

type PlayerHandler struct {
	Player PlayerInfo
	Router router.Router
}

func (ph *PlayerHandler) Init(ctx context.Context, request *player_pb.InitRequest) (*player_pb.InitResponse, error) {
	newBlobId, startX, startY, mass := ph.Player.NewBlob()
	response := player_pb.InitResponse{
		Id:   newBlobId,
		X:    startX,
		Y:    startY,
		Mass: mass,
	}
	
	return &response, nil
}

// Search function responsible to get the Country information
func (ph *PlayerHandler) Move(ctx context.Context, request *player_pb.MoveRequest) (*player_pb.MoveResponse, error) {
	// for now just echo response with increment on position
	// log.Println("Moving!")
	if !ph.Player.GetAlive() {
		response := player_pb.MoveResponse{
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

	response := player_pb.MoveResponse{
		X:     x,
		Y:     y,
		Alive: true,
		Mass:  ph.Player.GetMass(),//TODO: change but leave for now
	}

	return &response, nil
}

func (ph *PlayerHandler) Region(ctx context.Context, request *player_pb.RegionRequest) (*player_pb.RegionResponse, error) {
	// regionId := request.GetId()
	// conn := ph.router.Get(regionId)
	// regionClient := region_pb.NewRegionClient(conn)

	// regionRequest := &region_pb.IdRegionRequest{Id: regionId}
	// getRegionResponse, _ := regionClient.GetRegion(context.Background(), regionRequest)

	// response := region_pb.RegionResponse{
	// 	Blobs: getRegionResponse.Blobs,
	// 	Foods: getRegionResponse.Foods,
	// }

	response := player_pb.RegionResponse{
		Blobs: []*player_pb.Blob{ph.Player.GetBlob()},
		Foods: make([]*player_pb.Food, 0),
	}

	return &response, nil
}

func (ph *PlayerHandler) MassIncrement(ctx context.Context, request *player_pb.MassIncrementRequest) (*player_pb.MassIncrementResponse, error) {
	deltaMass := request.GetMassIncrease()
	ph.Player.IncrementMass(deltaMass)

	return nil, nil
}