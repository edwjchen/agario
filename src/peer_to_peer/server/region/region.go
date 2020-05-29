package region 

import (
	"golang.org/x/net/context"
	// "log"
)

type Region struct{}

func (Region) Ping(ctx context.Context, request *PingRequest) (*PingResponse, error) {
	response := PingResponse{}
	return &response, nil
}