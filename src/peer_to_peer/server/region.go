package region 

import (
	"peer_to_peer/server/region"
)

type Region struct{}

func (Region) Ping(ctx context.Context, request *region.PingRequest) (*blob.PingResponse, error) {
	response := region.PingResponse{}
	return &response, nil
}