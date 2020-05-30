package region

import (
	"peer_to_peer/server/player"
	"sync"
	"time"
	"math/rand"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"
	"peer_to_peer/server/region_pb"
)

type RegionInfo struct {
	FoodTree *quadtree.Quadtree
	PlayersIn     map[string]*player.PlayerInfo
	PlayersSeen   map[string]*player.PlayerInfo
	foodMux       sync.Mutex
	playerInMux   sync.Mutex
	playerSeenMux sync.Mutex
	x        uint16   
	y        uint16   
	xmin     float64 
	xmax     float64
	ymin     float64
	ymax     float64
	hash     uint32   
	Quit     chan bool
}

func (r *RegionInfo) NewRegion(x uint16, y uint16, hash uint32) {
	r.foodMux.Lock()
	r.playerInMux.Lock()
	r.playerSeenMux.Lock()
	r.FoodTree = quadtree.New(orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{player.MAP_WIDTH, player.MAP_HEIGHT}})
	r.PlayersIn = make(map[string]*player.PlayerInfo)
	r.PlayersSeen = make(map[string]*player.PlayerInfo)
	r.x = x
	r.y = y
	r.xmin = x * 500.0
	r.xmax = (x + 1) * 500.0 
	r.ymin = y * 500.0
	r.ymax = (y + 1) * 500.0
	r.hash = hash
	go func() {
		<-time.Ticker(time.Second)
		for {
			select {
			case <-r.Quit:
				return
			default:
				r.spawnFood()
			}
		}
	}()
	r.foodMux.Unlock()
	r.playerInMux.Unlock()
	r.playerSeenMux.Unlock()
}

func (r *RegionInfo) GetFood() []*Food {
	r.foodMux.Lock()
	defer r.foodMux.Unlock()
	bound := orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{player.MAP_WIDTH, player.MAP_HEIGHT}}
	foods := f.FoodTree.InBound([]orb.Pointer{}, bound)

	foodSlice := make([]*Food, len(foods))
	for idx, food := range foods {
		point := food.Point()
		foodSlice[idx] = &Food{X: point.X(), Y: point.Y()}
	}
	return foodSlice
}

func (r *RegionInfo) GetIn() map[string]*player.PlayerInfo {
	r.playerInMux.Lock()
	defer r.playerInMux.Unlock()
	copy := r.PlayersIn
	return copy
}

func (r *RegionInfo) GetSeen() map[string]*player.PlayerInfo {
	r.playerSeenMux.Lock()
	defer r.playerSeenMux.Unlock()
	copy := r.PlayersSeen
	return copy
}

func (r *RegionInfo) spawnFood() {
	r.foodMux.Lock()
	defer r.foodMux.Unlock()

	spawnRandNum := rand.Intn(player.MIN_FOOD_NUM)

	for i := 0; i < spawnRandNum; i++ {
		x := float64(rand.Intn(player.MAP_WIDTH))
    	y := float64(rand.Intn(player.MAP_HEIGHT))

		foodPoint := orb.Point{x, y}

		r.FoodTree.Add(foodPoint)
	}
}


func (r *RegionInfo) removeFood(foodPointer orb.Pointer) {
	r.foodMux.Lock()
	defer r.foodMux.Unlock()
	r.FoodTree.Remove(foodPointer, nil)
}

// Returns number of foods eaten by player
func (r *RegionInfo) GetNumFoodsEaten(player *player.Player) int32 {
	// Delegate to removeFood
	// Get rectangular bound around player

	r.mux.Lock()
	defer r.mux.Unlock()
	radius := float64(player.Mass / 2)
	playerBound := orb.Bound{Min: orb.Point{player.X - radius, player.Y - radius}, Max: orb.Point{player.X + radius, player.Y + radius}}

	foodSlice := r.FoodTree.InBound([]orb.Pointer{}, playerBound)
	for _, food := range foodSlice {
		r.removeFood(food)
	}
	// log.Println("Eating: ", foodSlice)

	// bound := orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{MAP_WIDTH, MAP_HEIGHT}}
	// log.Println(f.FoodTree.InBound([]orb.Pointer{}, bound))

	return int32(len(foodSlice))
}

func GetRegionID(x, y uint16) uint32 {
	hasher := fnv.New32a()
	hasher.Write([]byte(uint32(x) << 16 | uint32(y)))
	hash := uint32(hasher.Sum32())
	return hash
}

func (r *RegionInfo) BlobIsIn(blob *Blob) bool {
	return r.xmin < blob.X && blob.X < r.xmax && r.ymin < blob.Y && blob.Y < r.ymax
}

func (r *RegionInfo) WasEaten(blob *blob.Blob) (bool, *Blob) {
	r.playerSeenMux.Lock()
	defer r.playerSeenMux.Unlock()

	blobRadius := float64(blob.Mass / 2)
	playerBound := orb.Bound{
		Min: orb.Point{blob.X - blobRadius, blob.Y - blobRadius}, 
		Max: orb.Point{blob.X + blobRadius, blob.Y + blobRadius},
	}

	for name, playerSeen := range PlayersSeen {
		if name == blob.Name {
			continue
		}

		currBlob := playerSeen.GetBlob()
		currBlobRadius := float64(currBlob.Mass / 2)

		centerDistance := blobDistance(blob.X, blob.Y, currBlob.X, currBlob.Y)

		if blobRadius > (centerDistance + currBlobRadius + EAT_RADIUS_DELTA) {
			return false, currBlob
		}
	}

	return true, nil
}


func blobDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((y2 - y1) * (y2 - y1) + (x2 - x1) * (x2 - x1))
}