package main

import (
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"net"
	"log"
	"os"
	"peer_to_peer/server/router"
	"peer_to_peer/server/region"
	"peer_to_peer/server/player"
	"peer_to_peer/server/player_pb"
	"peer_to_peer/server/region_pb"
	"peer_to_peer/entryserver"
)

const ENTRYSERVERIP = "192.168.86.24:8080"

func main() {

	myAddr := os.Args[1]
	regionAddr := myAddr+":"+player.REGION_PORT
	playerAddr := myAddr+":"+player.PLAYER_PORT
	
    router := &router.Router{}
	regionGrpcServer := grpc.NewServer()
	var regionHandler region.RegionHandler = region.RegionHandler{Router: router}
	region_pb.RegisterRegionServer(regionGrpcServer, &regionHandler)
	regionListener, err := net.Listen("tcp", regionAddr)
	if err != nil {
		log.Fatalf("could not listen to "+regionAddr+" %v", err)
	}
	log.Println("RegionHandler starting...")
	go regionGrpcServer.Serve(regionListener)

	conn, _ := grpc.Dial(ENTRYSERVERIP, grpc.WithInsecure())
	client := entryserver.NewEntryServerClient(conn)
	joinRequest := &entryserver.JoinRequest{Ip: myAddr}
	client.Join(context.Background(), joinRequest)

	canStartRequest := &entryserver.CanStartRequest{}
	for {
		res, _ := client.CanStart(context.Background(), canStartRequest)
		if res.GetCanStart() {
			break
		}
	}
	
	router.Init([]string{"192.168.86.24:3001", "192.168.86.25:3001"}, regionAddr)
	regionHandler.Init()
	
	log.Println("PlayerHandler starting to process...")
	var playerHandler player.PlayerHandler = player.PlayerHandler{Router: router}
	playerHandler.Player.InitIP(playerAddr)
	// player.PlayerInfoStruct.InitIP(playerAddr)
	playerGrpcServer := grpc.NewServer()
	
	player_pb.RegisterPlayerServer(playerGrpcServer, &playerHandler)
	playerListener, err := net.Listen("tcp", playerAddr)
	if err != nil {
		log.Fatalf("could not listen to "+playerAddr+" %v", err)
	}
	log.Println("PlayerHandler starting...")
	// go playerGrpcServer.Serve(playerListener)
	log.Fatal(playerGrpcServer.Serve(playerListener))


}