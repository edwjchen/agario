package region

import (
	"peer_to_peer/server/player"
	"sync"
	"math/rand"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"
)

type RegionInfo struct {
	foodTree *quadtree.Quadtree
	playersBelongs  map[string]*player.Player
	playersPresents map[string]*player.Player
	mux      sync.Mutex
	x        uint16   
	y        uint16   
	hash     uint32   
}

func (r *RegionInfo) NewRegion(x uint16, y uint16, hash uint32) {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.foodTree = quadtree.New(orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{player.MAP_WIDTH, player.MAP_HEIGHT}})
	r.playersBelongs = make(map[string]*player.Player)
	r.playersPresents = make(map[string]*player.Player)
	r.x = x
	r.y = y
	r.hash = hash
}

func (r *RegionInfo) GetFood() []*Food {
	r.mux.Lock()
	defer r.mux.Unlock()
	bound := orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{player.MAP_WIDTH, player.MAP_HEIGHT}}
	foods := f.foodTree.InBound([]orb.Pointer{}, bound)

	foodSlice := make([]*Food, len(foods))
	for idx, food := range foods {
		point := food.Point()
		foodSlice[idx] = &Food{X: point.X(), Y: point.Y()}
	}
	return foodSlice
}


func (r *RegionInfo) SpawnFood() {
	r.mux.Lock()
	defer r.mux.Unlock()

	spawnRandNum := rand.Intn(player.MIN_FOOD_NUM)

	for i := 0; i < spawnRandNum; i++ {
		x := float64(rand.Intn(player.MAP_WIDTH))
    	y := float64(rand.Intn(player.MAP_HEIGHT))

		foodPoint := orb.Point{x, y}

		r.foodTree.Add(foodPoint)
	}
}


func (r *RegionInfo) removeFood(foodPointer orb.Pointer) {
	r.foodTree.Remove(foodPointer, nil)
}

// Returns number of foods eaten by player
func (r *RegionInfo) GetNumFoodsEaten(player *player.Player) int32 {
	// Delegate to removeFood
	// Get rectangular bound around player

	r.mux.Lock()
	defer r.mux.Unlock()
	radius := float64(player.Mass / 2)
	playerBound := orb.Bound{Min: orb.Point{player.X - radius, player.Y - radius}, Max: orb.Point{player.X + radius, player.Y + radius}}

	foodSlice := r.foodTree.InBound([]orb.Pointer{}, playerBound)
	for _, food := range foodSlice {
		r.removeFood(food)
	}
	// log.Println("Eating: ", foodSlice)

	// bound := orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{MAP_WIDTH, MAP_HEIGHT}}
	// log.Println(f.foodTree.InBound([]orb.Pointer{}, bound))

	return int32(len(foodSlice))
}