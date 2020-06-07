package main

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"peer_to_peer/common"
	. "peer_to_peer/entryserver"
	"peer_to_peer/info"
	// "math"
)

const JOIN_ACTION = "JOIN"
const CREATE_ACTION = "CREATE"

var esConfig common.EntryServerConfig
var entryInfo info.EntryInfo

func main() {
	grpcServer := grpc.NewServer()
	var server EntryServer
	RegisterEntryServerServer(grpcServer, server)

	conf, err := common.ReadEntryServerConfig(os.Getenv("GOPATH") + "/src/peer_to_peer/common/esconfig.json")
	if err != nil {
		log.Fatalf("Can't read Entry server config", err)
	}
	esConfig = conf
	// esConfig.MAP_LENGTH = int32(math.Sqrt(conf.MAX_PLAYERS))
	entryInfo = info.EntryInfo{
		MinPlayers: int(esConfig.MIN_PLAYERS),
		MaxPlayers: int(esConfig.MAX_PLAYERS),
		CurrNodes:  make([]string, 0),
	}

	listen, err := net.Listen("tcp", esConfig.ADDR)
	if err != nil {
		log.Fatalf("could not listen to ", esConfig.ADDR, err)
	}
	log.Println("Server starting...")
	log.Fatal(grpcServer.Serve(listen))
}

type EntryServer struct{}

func (EntryServer) CanStart(ctx context.Context, request *CanStartRequest) (*CanStartReply, error) {
	/*
		if
	*/
	response := &CanStartReply{
		CanStart: entryInfo.CanStart(),
	}

	return response, nil

}

func (EntryServer) Join(ctx context.Context, request *JoinRequest) (*JoinReply, error) {
	//get number of players
	//if current number of players < MAX-players, return true
	//get current number of players
	//if == 0
	//else
	joinIp := request.GetIp()
	var action string
	var ip string
	if entryInfo.ShouldCreate() {
		action = JOIN_ACTION
		entryInfo.AddNode(joinIp)
		ip = entryInfo.GetNodeIP() //get IP CALL HERE
	} else {
		action = CREATE_ACTION
		ip = joinIp
		entryInfo.AddNode(joinIp)
	}
	response := &JoinReply{
		Action:    action,
		Ip:        ip,
		MapLength: 0,
	}
	return response, nil
}
