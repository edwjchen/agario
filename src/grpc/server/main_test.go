package main
import (
  "grpc/server/blob"
  "log"
  "testing"
  "golang.org/x/net/context"
)
func TestBlob(t *testing.T) {
  ctx := context.Background()
  request := blob.BlobRequest{&blob.Position{X:0, Y:0}}
  server := Server{}
  response, err := server.Move(ctx, &request)
  if err != nil {
    t.Error(err)
  }
  log.Println(response)
}