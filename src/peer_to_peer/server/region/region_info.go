package region

import (
	"peer_to_peer/server/player"
	"sync"
	"time"
	"math"
	"math/rand"
	"log"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"
	. "peer_to_peer/common"
)

type RegionInfo struct {
	FoodTree *quadtree.Quadtree
	PlayersIn     map[string]*player.PlayerInfo
	PlayersSeen   map[string]*player.PlayerInfo
	foodMux       sync.Mutex
	PlayerInMux   sync.Mutex
	PlayerSeenMux sync.Mutex
	x        uint16   
	y        uint16   
	xmin     float64 
	xmax     float64
	ymin     float64
	ymax     float64
	hash     uint32   
	Quit     chan bool
}

func (r *RegionInfo) InitRegion(x, y uint32) {
	log.Println("initing region for x:", x, " y:", y)
	r.foodMux.Lock()
	r.PlayerInMux.Lock()
	r.PlayerSeenMux.Lock()
	r.PlayersIn = make(map[string]*player.PlayerInfo)
	r.PlayersSeen = make(map[string]*player.PlayerInfo)
	r.x = uint16(x)
	r.y = uint16(y)
	r.xmin = float64(x * Conf.REGION_MAP_WIDTH)
	r.xmax = float64((x + 1) * Conf.REGION_MAP_WIDTH)
	r.ymin = float64(y * Conf.REGION_MAP_HEIGHT)
	r.ymax = float64((y + 1) * Conf.REGION_MAP_HEIGHT)
	r.FoodTree = quadtree.New(orb.Bound{Min: orb.Point{r.xmin, r.ymin}, Max: orb.Point{r.xmax, r.ymax}})
	r.foodMux.Unlock()
	r.PlayerInMux.Unlock()
	r.PlayerSeenMux.Unlock()
}

func (r *RegionInfo) RunSpawnFood() {
	// go func() {
	for {
		<-time.Tick(time.Second)
		select {
		case <-r.Quit:
			return
		default:
			r.spawnFood()
		}
	}
	// }()
}

func (r *RegionInfo) GetFood() []*Food {
	r.foodMux.Lock()
	defer r.foodMux.Unlock()
	bound := orb.Bound{Min: orb.Point{r.xmin, r.ymin}, Max: orb.Point{r.xmax, r.ymax}}
	foods := r.FoodTree.InBound([]orb.Pointer{}, bound)

	foodSlice := make([]*Food, len(foods))
	for idx, food := range foods {
		point := food.Point()
		foodSlice[idx] = &Food{X: point.X(), Y: point.Y()}
	}
	return foodSlice
}

func (r *RegionInfo) GetIn() map[string]*player.PlayerInfo {
	r.PlayerInMux.Lock()
	defer r.PlayerInMux.Unlock()
	copy := make(map[string]*player.PlayerInfo)
	for k,v := range r.PlayersIn {
		copy[k] = v
	}
	// copy := r.PlayersIn
	return copy
}

func (r *RegionInfo) GetSeen() map[string]*player.PlayerInfo {
	r.PlayerSeenMux.Lock()
	defer r.PlayerSeenMux.Unlock()
	copy := make(map[string]*player.PlayerInfo)
	for k,v := range r.PlayersSeen {
		copy[k] = v
	}
	// copy := r.PlayersSeen
	return copy
}

func (r *RegionInfo) spawnFood() {
	r.foodMux.Lock()
	defer r.foodMux.Unlock()

	// TODO: Check if we have enought food

	spawnRandNum := rand.Intn(2)

	for i := 0; i < spawnRandNum; i++ {
		x := float64(rand.Intn(int(Conf.REGION_MAP_WIDTH))) + r.xmin
		y := float64(rand.Intn(int(Conf.REGION_MAP_HEIGHT))) + r.ymin

		foodPoint := orb.Point{x, y}

		r.FoodTree.Add(foodPoint)
	}
}


func (r *RegionInfo) removeFood(foodPointer orb.Pointer) {
	r.FoodTree.Remove(foodPointer, nil)
}

// Returns number of foods eaten by player
func (r *RegionInfo) GetNumFoodsEaten(blob *Blob) int32 {
	// Delegate to removeFood
	// Get rectangular bound around blob

	r.foodMux.Lock()
	defer r.foodMux.Unlock()
	radius := float64(blob.Mass / 2)
	playerBound := orb.Bound{Min: orb.Point{blob.X - radius, blob.Y - radius}, Max: orb.Point{blob.X + radius, blob.Y + radius}}

	foodSlice := r.FoodTree.InBound([]orb.Pointer{}, playerBound)
	for _, food := range foodSlice {
		r.removeFood(food)
	}
	// log.Println("Eating: ", foodSlice)

	// bound := orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{MAP_WIDTH, MAP_HEIGHT}}
	// log.Println(f.FoodTree.InBound([]orb.Pointer{}, bound))

	return int32(len(foodSlice))
}

func getRegionID(x, y uint16) uint32 {
	return uint32(x) << 16 | uint32(y)
}

func (r *RegionInfo) BlobIsIn(blob *Blob) bool {
	return r.xmin < blob.X && blob.X < r.xmax && r.ymin < blob.Y && blob.Y < r.ymax
}

func (r *RegionInfo) WasEaten(blob *Blob) (bool, *Blob) {

	blobRadius := float64(blob.Mass / 2)
	// playerBound := orb.Bound{
	// 	Min: orb.Point{blob.X - blobRadius, blob.Y - blobRadius}, 
	// 	Max: orb.Point{blob.X + blobRadius, blob.Y + blobRadius},
	// }

	for ip, playerSeen := range r.PlayersSeen {
		if ip == blob.Ip {
			continue
		}

		currBlob := playerSeen.GetBlob()
		currBlobRadius := float64(currBlob.Mass / 2)

		centerDistance := blobDistance(blob.X, blob.Y, currBlob.X, currBlob.Y)

		if blobRadius > (centerDistance + currBlobRadius + Conf.EAT_RADIUS_DELTA) {
			return false, currBlob
		}
	}

	return true, nil
}


func blobDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((y2 - y1) * (y2 - y1) + (x2 - x1) * (x2 - x1))
}