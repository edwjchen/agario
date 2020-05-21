package info

import (
	"grpc/server/blob"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"
	"math/rand"
	"strconv"
	"sync"
	// "log"
)

const SCREEN_WIDTH = 800
const SCREEN_HEIGHT = 500
const STARTING_MASS = 20
const SERVER_ID = "server1::"

type BlobsInfo struct {
	blobs    map[string]*blob.Player
	blobTree *quadtree.Quadtree
	mux      sync.Mutex
}

func (b *BlobsInfo) InitBlobs() {
	// TODO change to map
	b.mux.Lock()
	b.blobs = make(map[string]*blob.Player)
	b.blobTree = quadtree.New(orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{SCREEN_WIDTH, SCREEN_HEIGHT}})
	b.mux.Unlock()
}

func (b *BlobsInfo) NewBlob() (string, float64, float64) {
	b.mux.Lock()
	defer b.mux.Unlock()
	newBlobId := SERVER_ID + strconv.Itoa(len(b.blobs))
	startX := rand.Float64()*400 + 100
	startY := rand.Float64()*400 + 100
	//add blob to map
	newBlob := &blob.Player{Id: newBlobId, X: startX, Y: startY, Alive: true, Mass: STARTING_MASS}
	b.blobs[newBlobId] = newBlob
	return newBlobId, startX, startY
}

func (b *BlobsInfo) UpdatePos(name string, dx float64, dy float64) (float64, float64) {
	b.mux.Lock()
	defer b.mux.Unlock()
	updateBlob := *b.blobs[name]
	updateBlob.X += dx
	updateBlob.Y += dy

	//constrain movement for now
	if updateBlob.X > SCREEN_WIDTH * 2 {
		updateBlob.X = SCREEN_WIDTH
	} else if updateBlob.X < 0 {
		updateBlob.X = 0
	}

	if updateBlob.Y > SCREEN_HEIGHT * 2 {
		updateBlob.Y = SCREEN_HEIGHT
	} else if updateBlob.Y < 0 {
		updateBlob.Y = 0
	}
	// log.Println(name, " is at pos ", updateBlob.x, updateBlob.y)

	*b.blobs[name] = updateBlob
	return updateBlob.X, updateBlob.Y
}

func (b *BlobsInfo) GetBlobs() []*blob.Player {
	b.mux.Lock()
	defer b.mux.Unlock()
	retBlobs := make([]*blob.Player, 0)
	// log.Println("Printing blobs")
	for _, blob := range b.blobs {
		// log.Println(blob)
		retBlobs = append(retBlobs, blob)
	}
	// log.Println("list of blobs", retBlobs)
	return retBlobs
}

// Returns update mass of blob
func (b *BlobsInfo) UpdateBlobMass(id string, foodInfo FoodInfo) int32 {
	b.mux.Lock()
	defer b.mux.Unlock()
	player := b.blobs[id]
	oldMass := player.Mass
	newMass := oldMass + foodInfo.GetNumFoodsEaten(player)
	b.blobs[id].Mass = newMass

	return newMass
}
