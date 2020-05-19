package main
import (
  "encoding/json"
  "grpc/server/countries"
  "io/ioutil"
  "log"
  "net"
  "net/http"
  "golang.org/x/net/context"
  "google.golang.org/grpc"
)
func main() {
  grpcServer := grpc.NewServer()
  var server Server
  countries.RegisterCountryServer(grpcServer, server)
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
func (Server) Search(ctx context.Context, request *countries.CountryRequest) (*countries.CountryResponse, error) {
  resp, err := http.Get("https://restcountries.eu/rest/v2/name/ " + request.Name)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  jsonData, err := ioutil.ReadAll(resp.Body)
  if err != nil {
    return nil, err
  }
  var data []countries.CountryResponse
  if err := json.Unmarshal(jsonData, &data); err != nil {
    return nil, err
  }
  return &data[0], nil
}