package main

import (
	"google.golang.org/grpc"
	"net"
	"log"
)

func main() {
	grpcServer := grpc.NewServer()
	listen, err := net.Listen("tcp", "0.0.0.0:3000")
	if err != nil {
		log.Fatalf("could not listen to 0.0.0.0:3000 %v", err)
	}
  	log.Println("Server starting...")
  	log.Fatal(grpcServer.Serve(listen))
}

type EntryServer struct {}

func (EntryServer) Init(ctx context.Context, request *blob.InitRequest) (*blob.InitResponse, error) {
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