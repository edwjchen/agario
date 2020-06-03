package main

import (
	"google.golang.org/grpc"
	"net"
	"log"
	"golang.org/x/net/context"
	"peer_to_peer/info"
	. "peer_to_peer/entryserver"
	"math"
)

const JOIN_ACTION = "JOIN"
const CREATE_ACTION = "CREATE"

const MIN_PLAYERS = 2
const MAX_PLAYERS = 10000
var MAP_LENGTH int32 = int32(math.Sqrt(MAX_PLAYERS))

var entryInfo info.EntryInfo = info.EntryInfo{MinPlayers: MIN_PLAYERS, MaxPlayers: MAX_PLAYERS, CurrNodes: make([]string, 0)}

func main() {
	grpcServer := grpc.NewServer()
	var server EntryServer
	RegisterEntryServerServer(grpcServer, server)
	
	listen, err := net.Listen("tcp", "0.0.0.0:8080")
	if err != nil {
		log.Fatalf("could not listen to localhost:8080 %v", err)
	}
  	log.Println("Server starting...")
  	log.Fatal(grpcServer.Serve(listen))
}

type EntryServer struct{}

func (EntryServer) CanStart(ctx context.Context, request *CanStartRequest) (*CanStartReply, error) {
	/*
	if 
	*/
	response := &CanStartReply {
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
		Action: action,
		Ip: ip,
		MapLength: MAP_LENGTH,
	}
	return response, nil
} 
