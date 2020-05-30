package player

import (
	"golang.org/x/net/context"
	"math"
)

type PlayerHandler struct {

}

const SCREEN_WIDTH = 800
const SCREEN_HEIGHT = 500

const speed = 4

var PlayerInfoStruct PlayerInfo

func (PlayerHandler) Init(ctx context.Context, request *InitRequest) (*InitResponse, error) {
	newBlobId, startX, startY, mass := PlayerInfoStruct.NewBlob()
	response := InitResponse{
		Id:   newBlobId,
		X:    startX,
		Y:    startY,
		Mass: mass,
	}
	
	return &response, nil
}

// Search function responsible to get the Country information
func (PlayerHandler) Move(ctx context.Context, request *MoveRequest) (*MoveResponse, error) {
	//for now just echo response with increment on position

	if !PlayerInfoStruct.GetAlive() {
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

	// log.Println("get: ", dx, dy)
	rotation := math.Atan2(dy-SCREEN_HEIGHT/2, dx-SCREEN_WIDTH/2) * 180 / math.Pi
	vx := speed * (90 - math.Abs(rotation)) / 90
	var vy float64
	if rotation < 0 {
		vy = -1*speed + math.Abs(vx)
	} else {
		vy = speed - math.Abs(vx)
	}

	x, y := PlayerInfoStruct.UpdatePos(vx, vy)

	response := MoveResponse{
		X:     x,
		Y:     y,
		Alive: true,
		Mass:  PlayerInfoStruct.GetMass(),//TODO: change but leave for now
	}

	return &response, nil
}

func (PlayerHandler) Region(ctx context.Context, request *RegionRequest) (*RegionResponse, error) {

	response := RegionResponse{
		Blobs: make([]*Blob, 0),
		Foods: make([]*Food, 0),
	}

	return &response, nil
}

func (PlayerHandler) MassIncrement(ctx context.Context, request *MassIncrementRequest) (*MassIncrementResponse, error) {
	return nil, nil
}