package main

import (
	"google.golang.org/grpc"
	"net"
	"log"
	"os"
	"peer_to_peer/server/router"
	"peer_to_peer/server/region"
	"peer_to_peer/server/player"
	"peer_to_peer/server/player_pb"
	"peer_to_peer/server/region_pb"
)

func main() {

	playerAddr := os.Args[1]
	
	router := &router.Router{}
	router.Init([]string{"localhost:3001"})
	go router.Heartbeat()
	var playerHandler player.PlayerHandler = player.PlayerHandler{Router: router}
	playerHandler.Player.InitIP(playerAddr)
	// player.PlayerInfoStruct.InitIP(playerAddr)
	playerGrpcServer := grpc.NewServer()
	
	player_pb.RegisterPlayerServer(playerGrpcServer, &playerHandler)
	playerListener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatalf("could not listen to 0.0.0.0:3000 %v", err)
	}
	log.Println("PlayerHandler starting...")
	go playerGrpcServer.Serve(playerListener)
	// log.Fatal(playerGrpcServer.Serve(playerListener))

	regionGrpcServer := grpc.NewServer()
	var regionHandler region.RegionHandler = region.RegionHandler{Router: router}
	region_pb.RegisterRegionServer(regionGrpcServer, &regionHandler)
	regionListener, err := net.Listen("tcp", "0.0.0.0:3001")
	if err != nil {
		log.Fatalf("could not listen to 0.0.0.0:3001 %v", err)
	}
	log.Println("RegionHandler starting...")
	log.Fatal(regionGrpcServer.Serve(regionListener))

}