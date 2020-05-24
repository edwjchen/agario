package info

import (
	"client_server/server/blob"
	"github.com/paulmach/orb"
	"github.com/paulmach/orb/quadtree"
	"math/rand"
	"sync"
	// "log"
)

const MIN_FOOD_NUM = 50

type FoodInfo struct {
	foodTree *quadtree.Quadtree
	// foodMap  map[blob.Food]*blob.Food
	mux      sync.Mutex
}

func (f *FoodInfo) InitFood() {
	// TODO change to map
	f.mux.Lock()
	// f.foodMap = make(map[blob.Food]*blob.Food)
	f.foodTree = quadtree.New(orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{MAP_WIDTH, MAP_HEIGHT}})
	f.mux.Unlock()

	f.SpawnFood()
}

// Doesn't spawn food if not needed
func (f *FoodInfo) SpawnFood() {
	f.mux.Lock()
	defer f.mux.Unlock()
	// bound := orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{MAP_WIDTH, MAP_HEIGHT}}
	// foods := f.foodTree.InBound([]orb.Pointer{}, bound)

	// if len(foods) > MIN_FOOD_NUM {
	// 	return
	// }
	// rand.Seed(0)

	spawnRandNum := rand.Intn(MIN_FOOD_NUM)

	for i := 0; i < spawnRandNum; i++ {
		x := float64(rand.Intn(MAP_WIDTH))
    	y := float64(rand.Intn(MAP_HEIGHT))

		// food := blob.Food{X: x, Y: y}
		// _, exists := f.foodMap[food]
		// if exists {
		// 	continue
		// }

		foodPoint := orb.Point{x, y}

		f.foodTree.Add(foodPoint)
		// f.foodMap[food] = &food
	}
}

func (f *FoodInfo) removeFood(foodPointer orb.Pointer) {
	f.foodTree.Remove(foodPointer, nil)

	// foodPoint := foodPointer.Point()
	// food := blob.Food{X: foodPoint.X(), Y: foodPoint.Y()}
	// log.Println(food)
	// delete(f.foodMap, food)

	// log.Println(len(f.foodMap))
}

// Returns number of foods eaten by player
func (f *FoodInfo) GetNumFoodsEaten(player *blob.Player) int32 {
	// Delegate to removeFood
	// Get rectangular bound around player

	f.mux.Lock()
	defer f.mux.Unlock()
	radius := float64(player.Mass / 2)
	playerBound := orb.Bound{Min: orb.Point{player.X - radius, player.Y - radius}, Max: orb.Point{player.X + radius, player.Y + radius}}

	foodSlice := f.foodTree.InBound([]orb.Pointer{}, playerBound)
	for _, food := range foodSlice {
		f.removeFood(food)
	}
	// log.Println("Eating: ", foodSlice)

	// bound := orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{MAP_WIDTH, MAP_HEIGHT}}
	// log.Println(f.foodTree.InBound([]orb.Pointer{}, bound))

	return int32(len(foodSlice))
}

func (f *FoodInfo) GetFoods() []*blob.Food {
	f.mux.Lock()
	defer f.mux.Unlock()

	bound := orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{MAP_WIDTH, MAP_HEIGHT}}
	foods := f.foodTree.InBound([]orb.Pointer{}, bound)

	foodSlice := make([]*blob.Food, len(foods))
	for idx, food := range foods {
		point := food.Point()
		foodSlice[idx] = &blob.Food{X: point.X(), Y: point.Y()}
	}

	return foodSlice
}
