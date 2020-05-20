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
// Search function responsible to get the Country information
func (Server) Move(ctx context.Context, request *blob.BlobRequest) (*blob.BlobResponse, error) {
  //for now just echo response with increment on position
  x := request.Position.GetX()
  y := request.Position.GetY()
  x++
  y++
  newPos := blob.Position{X: x, Y: y}
  response := blob.BlobResponse{
    Position: &newPos,
    Alive: true,
    Mass: 0,
    Players: make([]byte, 0),
    Food: make([]byte, 0)}

  return &response, nil
}