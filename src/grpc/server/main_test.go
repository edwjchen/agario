package main
import (
  "grpc/server/countries"
  "log"
  "testing"
  "golang.org/x/net/context"
)
func TestCountry(t *testing.T) {
  ctx := context.Background()
  request := countries.CountryRequest{Name: "Brazil"}
  server := Server{}
  response, err := server.Search(ctx, &request)
  if err != nil {
    t.Error(err)
  }
  if response.Alpha2Code != "BR" {
    t.Error("Different Country returned")
  }
  log.Println(response)
}