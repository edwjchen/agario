package info

import (
	"grpc/server/blob"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"
	"math/rand"
	"sync"
	"log"
)

const MIN_FOOD_NUM = 50

type FoodInfo struct {
	foodTree *quadtree.Quadtree
	foodMap  map[blob.Food]*blob.Food
	mux      sync.Mutex
}

func (f *FoodInfo) InitFood() {
	// TODO change to map
	f.mux.Lock()
	f.foodMap = make(map[blob.Food]*blob.Food)
	f.foodTree = quadtree.New(orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{SCREEN_WIDTH, SCREEN_HEIGHT}})
	f.mux.Unlock()

	f.SpawnFood()
}

// Doesn't spawn food if not needed
func (f *FoodInfo) SpawnFood() {
	f.mux.Lock()
	defer f.mux.Unlock()
	if len(f.foodMap) > MIN_FOOD_NUM {
		return
	}

	spawnRandNum := rand.Intn(MIN_FOOD_NUM-len(f.foodMap))

	for i := 0; i < spawnRandNum; i++ {
		x := rand.Float64() * SCREEN_WIDTH
    	y := rand.Float64() * SCREEN_HEIGHT

		foodPoint := orb.Point{x, y}
		food := blob.Food{X: x, Y: y}

		f.foodTree.Add(foodPoint)
		f.foodMap[food] = &food
	}
}

func (f *FoodInfo) removeFood(foodPointer orb.Pointer) {
	f.foodTree.Remove(foodPointer, nil)

	foodPoint := foodPointer.Point()
	food := blob.Food{X: foodPoint.X(), Y: foodPoint.Y()}
	log.Println(food)
	delete(f.foodMap, food)

	log.Println(len(f.foodMap))
}

// Returns number of foods eaten by player
func (f *FoodInfo) GetNumFoodsEaten(player *blob.Player) int32 {
	// Delegate to removeFood
	// Get rectangular bound around player


	log.Println("waiting on f lock")
	f.mux.Lock()
	defer f.mux.Unlock()
	log.Println("i got the f lock")
	radius := float64(player.Mass / 2)
	playerBound := orb.Bound{Min: orb.Point{player.X - radius, player.Y - radius}, Max: orb.Point{player.X + radius, player.Y + radius}}

	foodSlice := f.foodTree.InBound([]orb.Pointer{}, playerBound)
	for _, food := range foodSlice {
		f.removeFood(food)
	}

	return int32(len(foodSlice))
}

func (f *FoodInfo) GetFoods() []*blob.Food {
	f.mux.Lock()
	defer f.mux.Unlock()

	foodSlice := make([]*blob.Food, len(f.foodMap))
	idx := 0
	for _, ptr := range f.foodMap {
		foodSlice[idx] = ptr
		idx++
	}
  
	return foodSlice
}
