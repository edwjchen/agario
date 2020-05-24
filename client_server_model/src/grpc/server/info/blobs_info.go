package info

import (
	"grpc/server/blob"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"
	"math/rand"
	"strconv"
	"sync"
	// "log"
	"math"
)

const SCREEN_WIDTH = 10000
const SCREEN_HEIGHT = 10000
const STARTING_MASS = 20
const EAT_RADIUS_DELTA = 5
const SERVER_ID = "server1::"

type BlobsInfo struct {
	blobsMap map[string]*blob.Player
	
	// pointsMap point of player mapped to playerPointer
	pointsMap map[orb.Point]*blob.Player

	blobTree *quadtree.Quadtree
	mux sync.Mutex
}

func (b *BlobsInfo) InitBlobs() {
	// TODO change to map
	b.mux.Lock()
	b.blobsMap = make(map[string]*blob.Player)
	b.blobTree = quadtree.New(orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{SCREEN_WIDTH, SCREEN_HEIGHT}})
	b.pointsMap = make(map[orb.Point]*blob.Player)
	b.mux.Unlock()
}

func (b *BlobsInfo) NewBlob() (string, float64, float64) {
	b.mux.Lock()
	defer b.mux.Unlock()
	newBlobId := SERVER_ID + strconv.Itoa(len(b.blobsMap))
	startX := rand.Float64()*400 + 100
	startY := rand.Float64()*400 + 100
	//add blob to map
	newBlob := &blob.Player{Id: newBlobId, X: startX, Y: startY, Alive: true, Mass: STARTING_MASS}
	b.blobsMap[newBlobId] = newBlob
	point := orb.Point{startX, startY}
	b.pointsMap[point] = newBlob
	b.blobTree.Add(point)

	return newBlobId, startX, startY
}

func (b *BlobsInfo) UpdatePos(name string, dx float64, dy float64) (float64, float64) {
	b.mux.Lock()
	defer b.mux.Unlock()
	updateBlob := b.blobsMap[name]

	oldPoint := orb.Point{updateBlob.X, updateBlob.Y}
	delete(b.pointsMap, oldPoint)

	updateBlob.X += dx
	updateBlob.Y += dy

	//constrain movement for now
	if updateBlob.X > SCREEN_WIDTH {
		updateBlob.X = SCREEN_WIDTH
	} else if updateBlob.X < 0 {
		updateBlob.X = 0
	}

	if updateBlob.Y > SCREEN_HEIGHT {
		updateBlob.Y = SCREEN_HEIGHT
	} else if updateBlob.Y < 0 {
		updateBlob.Y = 0
	}
	// log.Println(name, " is at pos ", updateBlob.x, updateBlob.y)
	b.blobsMap[name] = updateBlob
	newPoint := orb.Point{updateBlob.X, updateBlob.Y}
	b.pointsMap[newPoint] = updateBlob

	b.blobTree.Remove(oldPoint, nil)
	b.blobTree.Add(newPoint)

	return updateBlob.X, updateBlob.Y
}

func (b *BlobsInfo) GetBlobs() []*blob.Player {
	b.mux.Lock()
	defer b.mux.Unlock()
	retBlobs := make([]*blob.Player, 0)
	// log.Println("Printing blobs")
	for _, blob := range b.blobsMap {
		if blob.Alive{
			// log.Println(blob)
			retBlobs = append(retBlobs, blob)
		}
	}
	// log.Println("list of blobs", retBlobs)
	return retBlobs
}

// Returns update mass of blob
func (b *BlobsInfo) UpdateBlobMass(id string, foodInfo *FoodInfo) int32 {
	b.mux.Lock()
	defer b.mux.Unlock()
	player := b.blobsMap[id]
	oldMass := player.Mass
	numFoodEaten := foodInfo.GetNumFoodsEaten(player)
	randNum := rand.Intn(100)
	if numFoodEaten > 0 && randNum == 0 && oldMass > STARTING_MASS + STARTING_MASS {
		player.Mass = oldMass - oldMass/2
	} else if numFoodEaten > 0 {
		player.Mass = oldMass + numFoodEaten
	}
	return player.Mass
	// log.Println("newMass: ", newMass)
}

func (b *BlobsInfo) EatBlobs(id string) {
	b.mux.Lock()
	defer b.mux.Unlock()

	currBlob := b.blobsMap[id]
	currBlobRadius := float64(currBlob.Mass / 2)
	playerBound := orb.Bound{Min: orb.Point{currBlob.X - currBlobRadius, currBlob.Y - currBlobRadius}, Max: orb.Point{currBlob.X + currBlobRadius, currBlob.Y + currBlobRadius}}

	blobsInBound := b.blobTree.InBound([]orb.Pointer{}, playerBound)
	for _, blobPoint := range blobsInBound {
		blob := b.pointsMap[blobPoint.Point()]
		blobRadius := float64(blob.Mass/2)

		centerDistance := blobDistance(currBlob.X, currBlob.Y, blob.X, blob.Y)

		if currBlobRadius > (centerDistance + blobRadius + EAT_RADIUS_DELTA) {
			currBlob.Mass += (blob.Mass) / 2
			b.removeBlob(blobPoint, blob.Id)
		}
	}
}

func (b *BlobsInfo) IsBlobAlive(id string) bool {
	b.mux.Lock()
	defer b.mux.Unlock()

	return b.blobsMap[id].Alive
}

func (b *BlobsInfo) removeBlob(blobPointer orb.Pointer, id string){ 
	delete(b.pointsMap, blobPointer.Point())
	b.blobTree.Remove(blobPointer, nil)
	b.blobsMap[id].Alive = false
}


func blobDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((y2 - y1) * (y2 - y1) + (x2 - x1) * (x2 - x1))
}