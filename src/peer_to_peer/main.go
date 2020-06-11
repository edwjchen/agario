package main

import (
	"log"
	"net"
	"os"
	"peer_to_peer/common"
	"peer_to_peer/entryserver"
	"peer_to_peer/server/player"
	"peer_to_peer/server/player_pb"
	"peer_to_peer/server/region"
	"peer_to_peer/server/region_pb"
	"peer_to_peer/server/router"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func main() {

	myAddr := os.Args[1]
	basePath := "/src/peer_to_peer/common/"
	configPath := "/src/peer_to_peer/common/config.json"
	if len(os.Args) == 3 {
		configPath = basePath + os.Args[2]
	}

	e := common.ReadConfig(os.Getenv("GOPATH") + configPath)
	if e != nil {
		log.Fatalf("Cannot load config file", e)
	}
	regionAddr := myAddr + ":" + common.Conf.REGION_PORT
	playerAddr := myAddr + ":" + common.Conf.PLAYER_PORT

	router := &router.Router{RegionChange: make(chan router.RegionChangeInfo)}
	regionGrpcServer := grpc.NewServer()
	regionHandler := region.RegionHandler{Router: router, RegionChange: router.RegionChange}
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
	joinResp, errabc := client.Join(context.Background(), joinRequest)
	log.Println(joinResp, errabc)

	canStartRequest := &entryserver.CanStartRequest{}
	for {
		res, _ := client.CanStart(context.Background(), canStartRequest)
		// log.Println(err)
		// log.Println(res)
		// log.Println(res.GetCanStart())
		if res.GetCanStart() {
			break
		}
	}
	router.Init(common.Conf.REGION_SERVERS, regionAddr, )

	if joinResp.GetAction() == "JOIN" {
		regionHandler.Join()
	} else {
		regionHandler.Init()
	}

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
