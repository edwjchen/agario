package info

import (
	"client_server/server/blob"
	"math/rand"
	"strconv"
	"sync"	
	"log"
	"math"
)

const SCREEN_WIDTH = 800
const SCREEN_HEIGHT = 500
const MAP_WIDTH = 1500
const MAP_HEIGHT = 1500
const REGION_MAP_WIDTH = 500
const REGION_MAP_HEIGHT = 500
const MAX_FOOD_NUM = 30
const STARTING_MASS = 20
const MASS_MULTIPLIER = 3.0
const EAT_RADIUS_DELTA = 5
const ZOOM = 100
const SERVER_ID = "server1::"

type BlobsInfo struct {
	// pointsMap point of player mapped to playerPointer
	regionMap map[Point]map[string]*blob.Player

	blobsMap map[string]*blob.Player

	mux sync.RWMutex
}

func (b *BlobsInfo) InitBlobs() {
	// TODO change to map
	b.mux.Lock()
	defer b.mux.Unlock()

	b.blobsMap = make(map[string]*blob.Player)

	b.regionMap = make(map[Point]map[string]*blob.Player)
	norm_map_width := MAP_WIDTH / REGION_MAP_WIDTH 
	norm_map_height := MAP_HEIGHT / REGION_MAP_HEIGHT 

	for x := 0; x < norm_map_width; x++ {
		for y := 0; y < norm_map_height; y++ {
			playerMap := make(map[string]*blob.Player)
			point := Point{X: float64(x), Y: float64(y)}
			b.regionMap[point] = playerMap
		}
	}

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

	startRegionPoint := Point{X: 0, Y: 0}
	b.regionMap[startRegionPoint][newBlobId] = newBlob
	
	return newBlobId, startX, startY
}

func (b *BlobsInfo) UpdatePos(name string, dx float64, dy float64) (float64, float64) {
	b.mux.Lock()
	defer b.mux.Unlock()

	currBlob := b.blobsMap[name]
	currBlobId := currBlob.Id
	currBlobRegionPoint := Point{X: float64(int(currBlob.X / REGION_MAP_WIDTH)), Y: float64(int(currBlob.Y / REGION_MAP_HEIGHT))}

	// delete old point from regions map 
	delete(b.regionMap[currBlobRegionPoint], currBlobId)

	currBlob.X += dx
	currBlob.Y += dy

	//constrain movement for now
	if currBlob.X > MAP_WIDTH-1 {
		currBlob.X = MAP_WIDTH-1
	} else if currBlob.X < 0 {
		currBlob.X = 0
	}

	if currBlob.Y > MAP_HEIGHT-1 {
		currBlob.Y = MAP_HEIGHT-1
	} else if currBlob.Y < 0 {
		currBlob.Y = 0
	}
	// log.Println(name, " is at pos ", currBlob.x, currBlob.y)
	b.blobsMap[name] = currBlob

	newCurrBlobRegionPoint := Point{X: float64(int(currBlob.X / REGION_MAP_WIDTH)), Y: float64(int(currBlob.Y / REGION_MAP_HEIGHT))}
	b.regionMap[newCurrBlobRegionPoint][currBlobId] = currBlob

	return currBlob.X, currBlob.Y
}

func (b *BlobsInfo) GetBlobs(blobId string, foodInfo *FoodInfo) []*blob.Player {
	b.mux.RLock()
	defer b.mux.RUnlock()
	retBlobs := make([]*blob.Player, 0)
	// log.Println("Printing blobs")

	currBlob := b.blobsMap[blobId]
	regionPoints := GetAOI(currBlob)

	for _, regionPoint := range regionPoints {
		blobMap := b.regionMap[regionPoint]
		for _, blob := range blobMap {
			retBlobs = append(retBlobs, blob)
		}
	}

	return retBlobs
}

func (b *BlobsInfo) GetPlayer(blobId string) *blob.Player {
	b.mux.RLock()
	defer b.mux.RUnlock()
	// log.Println("Printing blobs")

	currBlob := b.blobsMap[blobId]
	copy := &blob.Player{Id: currBlob.Id, X: currBlob.X, Y: currBlob.Y, Alive: currBlob.Alive, Mass: currBlob.Mass}

	return copy
}

// Returns update mass of blob
func (b *BlobsInfo) UpdateBlobMass(blobId string, foodInfo *FoodInfo) int32 {
	b.mux.Lock()
	defer b.mux.Unlock()
	player := b.blobsMap[blobId]
	oldMass := player.Mass
	numFoodEaten := foodInfo.GetNumFoodsEaten(player)
	randNum := rand.Intn(1000)
	if numFoodEaten > 0 && randNum == 0 && oldMass > STARTING_MASS + STARTING_MASS {
		player.Mass = oldMass - oldMass/2
	} else if numFoodEaten > 0 {
		player.Mass = oldMass + numFoodEaten
	}
	return player.Mass
	// log.Println("newMass: ", newMass)
}

func (b *BlobsInfo) EatBlobs(blobId string) {
	b.mux.Lock()
	defer b.mux.Unlock()

	currBlob := b.blobsMap[blobId]
	currBlobRadius := GetRadiusFromMass(currBlob.Mass)

	regionPoints := GetAOI(currBlob)

	for _, regionPoint := range regionPoints {
		blobMap := b.regionMap[regionPoint]
		for _, blob := range blobMap {
			if currBlob.Id == blob.Id {
				continue
			}

			blobRadius := GetRadiusFromMass(blob.Mass)
			centerDistance := blobDistance(currBlob.X, currBlob.Y, blob.X, blob.Y)

			if currBlobRadius > (centerDistance + blobRadius + EAT_RADIUS_DELTA) {
				currBlob.Mass += (blob.Mass) / 2
				b.blobsMap[blob.Id].Alive = false
				delete(blobMap, blob.Id)
			}
		}
	}
}

func (b *BlobsInfo) IsBlobAlive(id string) bool {
	b.mux.RLock()
	defer b.mux.RUnlock()

	return b.blobsMap[id].Alive
}

func (b *BlobsInfo) removeBlob(id string) { 
	log.Println("remove blob: ", id)
	//get current pos
	//get region pos
	//remove from regionMap
	currBlob := b.blobsMap[id]
	b.blobsMap[id].Alive = false
	regionPoint := Point{X: float64(int(currBlob.X / MAP_WIDTH)), Y: float64(int(currBlob.Y / MAP_HEIGHT))}
	delete(b.regionMap[regionPoint], id)
}

func blobDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((y2 - y1) * (y2 - y1) + (x2 - x1) * (x2 - x1))
}

func GetRadiusFromMass(mass int32) float64 {
	rad := math.Sqrt(float64(mass)) * float64(MASS_MULTIPLIER)
	return rad
}

// Returns a list of region ids
func GetAOI(player *blob.Player) []Point {
	// p.mux.Lock()
	// defer p.mux.Unlock()
	NREGION_WIDTH := MAP_WIDTH / REGION_MAP_WIDTH
	NREGION_HEIGHT := MAP_HEIGHT / REGION_MAP_HEIGHT
 
	diameter := GetRadiusFromMass(player.Mass) * 2.0
	zoom_factor := ZOOM/diameter + 0.3
	top_left_x := player.X - float64(SCREEN_WIDTH)/zoom_factor/2
	top_left_y := player.Y - float64(SCREEN_HEIGHT)/zoom_factor/2
	bot_right_x := player.X + float64(SCREEN_WIDTH)/zoom_factor/2
	bot_right_y := player.Y + float64(SCREEN_HEIGHT)/zoom_factor/2
	//figure out which of the regions are in AOI
	start_region_x := float64(max(0, int32(math.Floor(top_left_x/float64(REGION_MAP_WIDTH)-1))))
	start_region_y := float64(max(0, int32(math.Floor(top_left_y/float64(REGION_MAP_HEIGHT)-1))))
	end_region_x := float64(min(int32(math.Ceil(bot_right_x/float64(REGION_MAP_WIDTH)+1)), int32(NREGION_WIDTH-1)))
	end_region_y := float64(min(int32(math.Ceil(bot_right_y/float64(REGION_MAP_HEIGHT)+1)), int32(NREGION_HEIGHT-1)))

	regionPoints := make([]Point, 0)

	for x := start_region_x; x <= end_region_x; x++ {
		for y := start_region_y; y <= end_region_y; y++ {
			regionPoints = append(regionPoints, Point{X: x, Y: y})
		}
	}

	return regionPoints
}

func min(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
