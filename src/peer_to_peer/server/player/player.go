package player


type Player struct{}

// const SCREEN_WIDTH = 10000
// const SCREEN_HEIGHT = 10000
const SCREEN_WIDTH = 800
const SCREEN_HEIGHT = 500
const STARTING_MASS = 20 //change later

const speed = 4
var x float64 = 400
var y float64 = 250

func (Player) Init(ctx context.Context, request *InitRequest) (*InitResponse, error) {
	newBlobId, startX, startY := blobsInfo.NewBlob()
	log.Println(newBlobId, "has joined")
	response := blob.InitResponse{
		Id:   newBlobId,
		X:    startX,
		Y:    startY,
		Mass: STARTING_MASS,
	}
	return &response, nil
}

// Search function responsible to get the Country information
func (Player) Move(ctx context.Context, request *blob.MoveRequest) (*blob.MoveResponse, error) {
	//for now just echo response with increment on position
	blobId := request.GetId()

	if !blobsInfo.IsBlobAlive(blobId) {
		response := blob.MoveResponse{
			X:     0,
			Y:     0,
			Alive: false,
			Mass:  0,
		}

		return &response, nil
	}

	dx := request.GetX()
	dy := request.GetY()

	// log.Println("get: ", dx, dy)
	rotation := math.Atan2(dy-SCREEN_HEIGHT/2, dx-SCREEN_WIDTH/2) * 180 / math.Pi
	vx := speed * (90 - math.Abs(rotation)) / 90
	var vy float64
	if rotation < 0 {
		vy = -1*speed + math.Abs(vx)
	} else {
		vy = speed - math.Abs(vx)
	}


	response := blob.MoveResponse{
		X:     x,
		Y:     y,
		Alive: true,
		Mass:  newMass,
	}

	return &response, nil
}

func (Player) Region(ctx context.Context, request *blob.RegionRequest) (*blob.RegionResponse, error) {

	response := blob.RegionResponse{
		Blobs: make([]*Blob, 0),
		Foods: make([]*Food, 0),
	}

	return &response, nil
}