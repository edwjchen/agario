package info

import (
	"client_server/server/blob"
	"math/rand"
	"sync"
	// "log"
)

type Point struct {
	X float64
	Y float64
}

type FoodInfo struct {
	// foodTree *quadtree.Quadtree
	// foodMap  map[blob.Food]*blob.Food
	
	//map of regions, each region contains a map of food
	//region points will be the top left corner == 0,0 for 1st region
	regionMap map[Point]map[Point]bool 
	mux      sync.Mutex
}

func (f *FoodInfo) InitFood() {
	// TODO change to map
	f.mux.Lock()
	defer f.mux.Unlock()
	// f.foodMap = make(map[blob.Food]*blob.Food)
	//f.foodTree = quadtree.New(orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{MAP_WIDTH, MAP_HEIGHT}})
	
	f.regionMap = make(map[Point]map[Point]bool)

	norm_map_width := MAP_WIDTH / REGION_MAP_WIDTH
	norm_map_height := MAP_HEIGHT / REGION_MAP_HEIGHT
	for x := 0; x < norm_map_width; x++ {
		for y := 0; y < norm_map_height; y++ {
			foodMap := make(map[Point]bool)
			point := Point{X: float64(x), Y: float64(y)}
			f.regionMap[point] = foodMap
			f.mux.Unlock()
			f.SpawnFood(point)
			f.mux.Lock()
		}
	}
}

// Doesn't spawn food if not needed
func (f *FoodInfo) SpawnFood(point Point) {
	f.mux.Lock()
	defer f.mux.Unlock()
	// bound := orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{MAP_WIDTH, MAP_HEIGHT}}
	// foods := f.foodTree.InBound([]orb.Pointer{}, bound)

	// if len(foods) > MIN_FOOD_NUM {
	// 	return
	// }
	// rand.Seed(0)

	// spawnRandNum := rand.Intn(MIN_FOOD_NUM)

	// for i := 0; i < spawnRandNum; i++ {
	// 	x := float64(rand.Intn(MAP_WIDTH))
    // 	y := float64(rand.Intn(MAP_HEIGHT))

	// 	// food := blob.Food{X: x, Y: y}
	// 	// _, exists := f.foodMap[food]
	// 	// if exists {
	// 	// 	continue
	// 	// }

	// 	foodPoint := orb.Point{x, y}

	// 	f.foodTree.Add(foodPoint)
	// 	// f.foodMap[food] = &food
	// }

	if len(f.regionMap[point]) > MAX_FOOD_NUM {
		return
	}

	spawnRandNum := rand.Intn(int(MAX_FOOD_NUM))

	for i := 0; i < spawnRandNum; i++ {
		x := float64(rand.Intn(int(REGION_MAP_WIDTH))) + point.X * REGION_MAP_WIDTH
		y := float64(rand.Intn(int(REGION_MAP_HEIGHT))) + point.Y * REGION_MAP_HEIGHT

		foodPoint := Point{X: float64(x), Y: float64(y)}

		f.regionMap[point][foodPoint] = true 
	}
}

func (f *FoodInfo) removeFood(point Point, foodPoint Point) {
	delete(f.regionMap[point], foodPoint)
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
	radius := GetRadiusFromMass(player.Mass)
	regionPoints := GetAOI(player)

	//get area of interest 
	foodEaten := []*blob.Food{}
	for _, regionPoint := range regionPoints {
		foodTree := f.regionMap[regionPoint]
		for foodPoint, _ := range foodTree {
			if blobDistance(player.X, player.Y, foodPoint.X, foodPoint.Y) <= radius {
				f.removeFood(regionPoint, foodPoint)
				foodEaten = append(foodEaten, &blob.Food{X:foodPoint.X, Y:foodPoint.Y})
			}
		}
	}
	// log.Println("Eating: ", foodSlice)

	// bound := orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{MAP_WIDTH, MAP_HEIGHT}}
	// log.Println(f.foodTree.InBound([]orb.Pointer{}, bound))

	return int32(len(foodEaten))
}

func (f *FoodInfo) GetFoods(player *blob.Player) []*blob.Food {
	f.mux.Lock()
	defer f.mux.Unlock()

	regionPoints := GetAOI(player)
	foodSlice := make([]*blob.Food, 0)

	for _, regionPoint := range regionPoints {
		foodMap := f.regionMap[regionPoint]
		for foodPoint, _ := range foodMap {
			food := &blob.Food{X: foodPoint.X, Y: foodPoint.Y}
			foodSlice = append(foodSlice, food)
		}
	}

	return foodSlice
}

