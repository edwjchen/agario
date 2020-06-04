package main

import (
	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"net"
	"log"
	"os"
	"peer_to_peer/common"
	"peer_to_peer/server/router"
	"peer_to_peer/server/region"
	"peer_to_peer/server/player"
	"peer_to_peer/server/player_pb"
	"peer_to_peer/server/region_pb"
	"peer_to_peer/entryserver"
)

// const ENTRYSERVERIP = "99.76.231.26:8080"

func main() {

	myAddr := os.Args[1]
	e := common.ReadConfig(os.Getenv("GOPATH")+"/src/peer_to_peer/common/config.json")
	if e != nil {
		log.Fatalf("Cannot load config file", e)
	}
	regionAddr := myAddr + ":" + common.Conf.REGION_PORT
	playerAddr := myAddr + ":" + common.Conf.PLAYER_PORT

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

	conn, err := grpc.Dial(common.Conf.ENTRY_SERVER, grpc.WithInsecure())
	log.Println(err)
	client := entryserver.NewEntryServerClient(conn)
	joinRequest := &entryserver.JoinRequest{Ip: myAddr}
	res, errabc := client.Join(context.Background(), joinRequest)
	log.Println(res, errabc)

	canStartRequest := &entryserver.CanStartRequest{}
	for {
		res, _:= client.CanStart(context.Background(), canStartRequest)
		// log.Println(err)
		// log.Println(res)
		// log.Println(res.GetCanStart())
		if res.GetCanStart() {
			break
		}
	}
	
	router.Init(common.Conf.REGION_SERVERS, regionAddr)
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