package main

import (
	// "encoding/json"
	"grpc/server/blob"
	// "io/ioutil"
	"log"
	"net"
	// "net/http"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
  "math"
  "grpc/server/info"
  "time"
	// "github.com/golang/protobuf/proto"
)

const SCREEN_WIDTH = 800
const SCREEN_HEIGHT = 500
const STARTING_MASS = info.STARTING_MASS

var foodInfo info.FoodInfo
var blobsInfo info.BlobsInfo


func main() {
	grpcServer := grpc.NewServer()
	var server Server
  blob.RegisterBlobServer(grpcServer, server)
  blobsInfo.InitBlobs()
  foodInfo.InitFood()
	listen, err := net.Listen("tcp", "0.0.0.0:3000")
	if err != nil {
		log.Fatalf("could not listen to 0.0.0.0:3000 %v", err)
	}
  log.Println("Server starting...")
  go spawnFood()
	log.Fatal(grpcServer.Serve(listen))
}

// Server is implementation proto interface
type Server struct{}

const speed = 4

var x float64 = 400
var y float64 = 250

func (Server) Init(ctx context.Context, request *blob.InitRequest) (*blob.InitResponse, error) {
	newBlobId, startX, startY := blobsInfo.NewBlob()
	log.Println(newBlobId, "has joined")
	response := blob.InitResponse{
		Id:   newBlobId,
		X:    startX,
		Y:    startY,
		Mass: STARTING_MASS,
	}
	return &response, nil
}

// Search function responsible to get the Country information
func (Server) Move(ctx context.Context, request *blob.MoveRequest) (*blob.MoveResponse, error) {
	//for now just echo response with increment on position
	blobId := request.GetId()

	if !blobsInfo.IsBlobAlive(blobId) {
		response := blob.MoveResponse{
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

	log.Println("send: ", x+vx, y+vy)
	x, y := blobsInfo.UpdatePos(blobId, vx, vy)
	newMass := blobsInfo.UpdateBlobMass(blobId, &foodInfo)
	blobsInfo.EatBlobs(blobId)
	// add func that gets if blob is alive.

	response := blob.MoveResponse{
		X:     x,
		Y:     y,
		Alive: true,
		Mass:  newMass,
	}

	return &response, nil
}

func (Server) Region(ctx context.Context, request *blob.RegionRequest) (*blob.RegionResponse, error) {
	// x := request.GetX()
	// y := request.GetY()
	// log.Println("pos: ", x, y)

	players := blobsInfo.GetBlobs()
	foods := foodInfo.GetFoods()
	response := blob.RegionResponse{
		Players: players,
		Foods:   foods,
	}

	return &response, nil
}

func spawnFood() {
  ticker := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		foodInfo.SpawnFood()
	}
}