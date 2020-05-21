package main
import (
  // "encoding/json"
  "grpc/server/blob"
  // "io/ioutil"
  "log"
  "net"
  // "net/http"
  "math"
  "golang.org/x/net/context"
  "google.golang.org/grpc"
  "strconv"
  "sync"
  "math/rand"
  "github.com/paulmach/orb"
  "github.com/paulmach/orb/quadtree"
  // "github.com/golang/protobuf/proto"
)

const SCREEN_WIDTH = 800
const SCREEN_HEIGHT = 500

func main() {
  grpcServer := grpc.NewServer()
  var server Server
  blob.RegisterBlobServer(grpcServer, server)
  foodInfo.InitFood()
  listen, err := net.Listen("tcp", "0.0.0.0:3000")
  if err != nil {
    log.Fatalf("could not listen to 0.0.0.0:3000 %v", err)
  }
  log.Println("Server starting...")
  log.Fatal(grpcServer.Serve(listen))
}
// Server is implementation proto interface
type Server struct{}

const MIN_FOOD_NUM = 800

type FoodInfo struct {
  foodTree *quadtree.Quadtree
  foodMap map[*blob.Food] bool
  mux sync.Mutex
}

var foodInfo FoodInfo

func (f *FoodInfo) InitFood() {
  // TODO change to map
  f.mux.Lock()
  f.foodMap = make(map[*blob.Food]bool)
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

  for i := 0 ; i < MIN_FOOD_NUM - len(f.foodMap) ; i++ {
    x := rand.Float64() * SCREEN_WIDTH
    y := rand.Float64() * SCREEN_HEIGHT
    
    foodPoint := orb.Point{x, y}
    food := &blob.Food{X: x, Y: y}
    
    f.foodTree.Add(foodPoint) 
    f.foodMap[food] = true
  }
}

func (f *FoodInfo) removeFood(foodPointer orb.Pointer) {
  f.foodTree.Remove(foodPointer, nil)

  foodPoint := foodPointer.Point()
  food := &blob.Food{X: foodPoint.X(), Y: foodPoint.Y()}
  delete(f.foodMap, food)
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

  return int32(len(foodSlice))
}

func (f *FoodInfo) GetFoods() []*blob.Food {
  f.mux.Lock()
  defer f.mux.Unlock()

  foodSlice := make([]*blob.Food, len(f.foodMap))
  idx := 0
  for food := range f.foodMap {
    foodSlice[idx] = food
    idx++
  }

  return foodSlice
}

type BlobsInfo struct {
  blobs map[string]*blob.Player
  blobTree quadtree.Quadtree
	mux sync.Mutex
}

func (b *BlobsInfo) NewBlob() (string, float64, float64) {
  b.mux.Lock()
  defer b.mux.Unlock()
  newBlobId := SERVER_ID + strconv.Itoa(len(b.blobs))
  startX := rand.Float64() * 400 + 100
  startY := rand.Float64() * 400 + 100
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

  *b.blobs[name] = updateBlob
  return updateBlob.X, updateBlob.Y
}

func (b *BlobsInfo) GetBlobs() []*blob.Player {
  b.mux.Lock()
  defer b.mux.Unlock()
  retBlobs := make([]*blob.Player, 0)
  log.Println("Printing blobs")
  for _, blob := range b.blobs {
    log.Println(blob)
    retBlobs = append(retBlobs, blob)
  }
  log.Println("list of blobs", retBlobs)
  return retBlobs 
}

// Returns update mass of blob
func (b * BlobsInfo) UpdateBlobMass(id string) int32 {
  b.mux.Lock()
  defer b.mux.Unlock()

  player := b.blobs[id]
  oldMass := player.Mass
  newMass := oldMass + foodInfo.GetNumFoodsEaten(player)
  b.blobs[id].Mass = newMass

  return newMass
}

const STARTING_MASS = 20
const SERVER_ID = "server1::"
const speed = 4

var x float64 = 400
var y float64 = 250

var blobsInfo BlobsInfo = BlobsInfo{blobs: make(map[string]*blob.Player)} 

func (Server) Init(ctx context.Context, request *blob.InitRequest) (*blob.InitResponse, error) {
  newBlobId, startX, startY := blobsInfo.NewBlob()
  log.Println(newBlobId, "has joined")
  response := blob.InitResponse {
    Id: newBlobId,
    X: startX,
    Y: startY,
    Mass: STARTING_MASS,
  }
  return &response, nil
}

// Search function responsible to get the Country information
func (Server) Move(ctx context.Context, request *blob.MoveRequest) (*blob.MoveResponse, error) {
  //for now just echo response with increment on position
  blobId := request.GetId()
  dx := request.GetX()
  dy := request.GetY()

  // log.Println("get: ", dx, dy)
  rotation := math.Atan2(dy - SCREEN_HEIGHT / 2,dx - SCREEN_WIDTH / 2) * 180 / math.Pi
  vx := speed * (90 - math.Abs(rotation)) / 90
  var vy float64
  if rotation < 0 {
    vy = -1 * speed + math.Abs(vx)
  } else {
    vy = speed  - math.Abs(vx)
  }

  // log.Println("send: ", x+vx, y+vy)
  x, y := blobsInfo.UpdatePos(blobId, vx, vy)
  newMass := blobsInfo.UpdateBlobMass(blobId)

  response := blob.MoveResponse {
    X: x,
    Y: y,
    Alive: true,
    Mass: newMass,
  }

  return &response, nil
}

func (Server) Region(ctx context.Context, request *blob.RegionRequest) (*blob.RegionResponse, error) {
  // x := request.GetX()
  // y := request.GetY()
  // log.Println("pos: ", x, y)

  response := blob.RegionResponse {
    Players: blobsInfo.GetBlobs(),
    Foods: foodInfo.GetFoods(),
  }

  return &response, nil
}

