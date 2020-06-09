package info

import (
	"client_server/server/blob"
	// "github.com/paulmach/orb"
	// "github.com/paulmach/orb/quadtree"
	"math/rand"
	"sync"
	// "log"
)

const MIN_FOOD_NUM = 50

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
	// f.foodMap = make(map[blob.Food]*blob.Food)
	//f.foodTree = quadtree.New(orb.Bound{Min: orb.Point{0, 0}, Max: orb.Point{MAP_WIDTH, MAP_HEIGHT}})
	
	norm_map_width := MAP_WIDTH / REGION_MAP_WIDTH
	norm_map_height := MAP_HEIGHT / REGION_MAP_HEIGHT
	for x := 0; x < norm_map_width; x++ {
		for y := 0; y < norm_map_height; y++ {
			foodMap := make(map[Point]bool)
			point := Point{X: float64(x), Y: float64(y)}
			f.regionMap[point] = foodMap

			f.SpawnFood(point)
		}
	}
	f.mux.Unlock()
}

// Doesn't spawn food if not needed
func (f *FoodInfo) SpawnFood(point Point) {
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
		x := float64(rand.Intn(int(REGION_MAP_WIDTH))) + point.X
		y := float64(rand.Intn(int(REGION_MAP_HEIGHT))) + point.Y

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

func blobDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt((y2 - y1) * (y2 - y1) + (x2 - x1) * (x2 - x1))
}
