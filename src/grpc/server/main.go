package main
import (
  // "encoding/json"
  "grpc/server/blob"
  // "io/ioutil"
  "log"
  "net"
  // "net/http"
  "math"
  "golang.org/x/net/context"
  "google.golang.org/grpc"
)
func main() {
  grpcServer := grpc.NewServer()
  var server Server
  blob.RegisterBlobServer(grpcServer, server)
  listen, err := net.Listen("tcp", "0.0.0.0:3000")
  if err != nil {
    log.Fatalf("could not listen to 0.0.0.0:3000 %v", err)
  }
  log.Println("Server starting...")
  log.Fatal(grpcServer.Serve(listen))
}
// Server is implementation proto interface
type Server struct{}

const SCREEN_WIDTH = 800
const SCREEN_HEIGHT = 500
const speed = 4

// Search function responsible to get the Country information
func (Server) Move(ctx context.Context, request *blob.BlobRequest) (*blob.BlobResponse, error) {
  //for now just echo response with increment on position
  x := request.Position.GetX()
  y := request.Position.GetY()

  rotation := math.Atan2(dY - SCREEN_HEIGHT / 2,dX - SCREEN_HEIGHT / 2) * 180 / math.Pi
  vx := speed * (90 - math.Abs(rotation)) / 90
  var vy float64
  if rotation < 0 {
    vy = -1 * speed + math.Abs(vx)
  } else {
    vy = speed  - math.Abs(vx)
  }

  newPos := blob.Position{X: x + vx, Y: y + vy}
  response := blob.BlobResponse{
    Position: &newPos,
    Alive: true,
    Mass: 0,
    Players: make([]byte, 0),
    Food: make([]byte, 0),
  }

  return &response, nil
}