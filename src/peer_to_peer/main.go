package main

import (
	"google.golang.org/grpc"
	"net"
	"log"
	"os"
	"peer_to_peer/server/player"
	// "peer_to_peer/server/region"
	// "peer_to_peer/server"
)

func main() {

	playerAddr := os.Args[1]
	
	player.PlayerInfoStruct.InitIP(playerAddr)
	playerGrpcServer := grpc.NewServer()
	var playerHandler player.PlayerHandler = player.PlayerHandler{}
	player.RegisterPlayerServer(playerGrpcServer, &playerHandler)
	playerListener, err := net.Listen("tcp", "localhost:3000")
	if err != nil {
		log.Fatalf("could not listen to 0.0.0.0:3000 %v", err)
	}
	log.Println("PlayerHandler starting...")
	// go playerGrpcServer.Server(playerListener)
	log.Fatal(playerGrpcServer.Serve(playerListener))

	// router := server.Router{}
	// router.Init([]string{"localhost:3001"})
	// go router.Heartbeat()

	// regionGrpcServer := grpc.NewServer()
	// var regionHandler region.RegionHandler
	// region.RegisterRegionHandlerServer(regionGrpcServer, regionHandler)
	// regionListener, err := net.Listen("tcp", "0.0.0.0:3001")
	// if err != nil {
	// 	log.Fatalf("could not listen to 0.0.0.0:3001 %v", err)
	// }
	// log.Println("RegionHandler Server starting...")
	// log.Fatal(regionGrpcServer.Serve(regionListener))

}