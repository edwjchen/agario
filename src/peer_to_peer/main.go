package main

import (
	"google.golang.org/grpc"
	"net"
	"log"
	// "peer_to_peer/server/player"
	"peer_to_peer/server/region"
	"peer_to_peer/server"
)

func main() {
	// playerGrpcServer := grpc.NewServer()
	// var playerServer player.PlayerServer
	// player.RegisterPlayerServer(playerGrpcServer, playerServer)
	// playerListener, err := net.Listen("tcp", "0.0.0.0:3000")
	// if err != nil {
	// 	log.Fatalf("could not listen to 0.0.0.0:3000 %v", err)
	// }
	// log.Println("Player Server starting...")
	// go playerGrpcServer.Serve(playerListener)

	router := server.Router{}
	router.Init([]string{"localhost:3001"})
	go router.Heartbeat()

	regionGrpcServer := grpc.NewServer()
	var regionServer region.Region
	region.RegisterRegionServer(regionGrpcServer, regionServer)
	regionListener, err := net.Listen("tcp", "0.0.0.0:3001")
	if err != nil {
		log.Fatalf("could not listen to 0.0.0.0:3001 %v", err)
	}
	log.Println("Region Server starting...")
	log.Fatal(regionGrpcServer.Serve(regionListener))

}