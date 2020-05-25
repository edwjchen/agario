package main

import (
	"google.golang.org/grpc"
	"net"
	"log"
	"sync"
	"golang.org/x/net/context"
	"peer_to_peer/entryserver"
	"peer_to_peer/info"
	"math"
)

const JOIN_ACTION = "JOIN"
const CREATE_ACTION = "CREATE"

const MIN_PLAYERS = 10
const MAX_PLAYERS = 10000
const MAP_LENGTH = int32(math.Sqrt(MAX_PLAYERS))

func main() {
	grpcServer := grpc.NewServer()
	listen, err := net.Listen("tcp", "0.0.0.0:3000")
	if err != nil {
		log.Fatalf("could not listen to 0.0.0.0:3000 %v", err)
	}
  	log.Println("Server starting...")
  	log.Fatal(grpcServer.Serve(listen))
}

entryInfo := EntryInfo{MinPlayers: MIN_PLAYERS, MaxPlayers: MAX_PLAYERS, CurrNodes: make([]string, 0)}

func (EntryServer) CanStart(ctx context.Context, request *entryserver.CanStartRequest) (*entryserver.CanStartReply, error) {
	/*
	if 
	*/
	response := &entryserver.CanStartReply {
		CanStart: entryInfo.CanStart()
	}

	return respnse, nil
	
}

func (EntryServer) Join(ctx context.Context, request *entryserver.JoinRequest) (*entryserver.JoinReply, error) {
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
		ip = //get IP CALL HERE
	} else {
		action = CREATE_ACTION
		ip = joinIp
	}
	response := entryserver.JoinReply{
		Action: action,
		Ip: i,
		MapLength: MAP_LENGTH,
	}
} 