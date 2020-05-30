package player

import (
	"golang.org/x/net/context"
	"math"
	// "log"
)

type PlayerHandler struct {

}

const SCREEN_WIDTH = 800
const SCREEN_HEIGHT = 500
const MIN_FOOD_NUM = 50


const speed = 4

var PlayerInfoStruct PlayerInfo

func (ph *PlayerHandler) Init(ctx context.Context, request *InitRequest) (*InitResponse, error) {
	// log.Println("got init")
	newBlobId, startX, startY, mass := PlayerInfoStruct.NewBlob()
	// log.Println("made new blob")
	response := InitResponse{
		Id:   newBlobId,
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

	// log.Println("dw & dy", dx, dy)

	// log.Println("get: ", dx, dy)
	rotation := math.Atan2(dy-SCREEN_HEIGHT/2, dx-SCREEN_WIDTH/2) * 180 / math.Pi
	vx := speed * (90 - math.Abs(rotation)) / 90
	var vy float64
	if rotation < 0 {
		vy = -1*speed + math.Abs(vx)
	} else {
		vy = speed - math.Abs(vx)
	}
	// log.Println("vx & vy", vx, vy)

	x, y := PlayerInfoStruct.UpdatePos(vx, vy)

	response := MoveResponse{
		X:     x,
		Y:     y,
		Alive: true,
		Mass:  PlayerInfoStruct.GetMass(),//TODO: change but leave for now
	}

	return &response, nil
}

func (ph *PlayerHandler) Region(ctx context.Context, request *RegionRequest) (*RegionResponse, error) {
	response := RegionResponse{
		Blobs: []*Blob{PlayerInfoStruct.GetBlob()},
		Foods: make([]*Food, 0),
	}

	return &response, nil
}

func (ph *PlayerHandler) MassIncrement(ctx context.Context, request *MassIncrementRequest) (*MassIncrementResponse, error) {
	return nil, nil
}