package main
import (
  "grpc/server/blob"
  "log"
  "testing"
  "time"
  "golang.org/x/net/context"
)
func TestBlob(t *testing.T) {
  ctx := context.Background()
  request := blob.BlobRequest{&blob.Position{X:0, Y:0}}
  server := Server{}
  start := time.Now()
  response, err := server.Move(ctx, &request)
  end := time.Now()
  log.Println(end.Sub(start))
  if err != nil {
    t.Error(err)
  }
  log.Println(response)
}