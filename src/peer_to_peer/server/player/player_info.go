package player

import (
	"io/ioutil"
	"math"
	"math/rand"
	. "peer_to_peer/common"
	"strconv"
	// "sync"
	"time"
)

type PlayerInfo struct {
	Blob       Blob
	LastUpdate time.Time
	// mux        sync.Mutex
	addr       string
}

func (p *PlayerInfo) InitIP(addr string) {
	// p.mux.Lock()
	// defer p.mux.Unlock()
	p.addr = addr
}

func (p *PlayerInfo) NewBlob() (string, float64, float64, int32, int32) {
	// TODO change to map
	// p.mux.Lock()
	// defer p.mux.Unlock()

	var x float64 = float64(rand.Intn(9999))
	var y float64 = float64(rand.Intn(9999))

	var ver int
	verBytes, err := ioutil.ReadFile(Conf.VER_FILE)
	if err != nil {
		ver = 0
	} else {
		ver, _ = strconv.Atoi(string(verBytes))
		ver++
	}

	writeVer := []byte(strconv.Itoa(ver))
	ioutil.WriteFile(Conf.VER_FILE, writeVer, 0666)

	p.Blob = Blob{Ip: p.addr, X: x, Y: y, Mass: Conf.STARTING_MASS, Alive: true, Ver: int32(ver), Seq: 0}
	return p.addr, x, y, Conf.STARTING_MASS, int32(ver)
}

func (p *PlayerInfo) GetAlive() bool {
	// p.mux.Lock()
	// defer p.mux.Unlock()
	return p.Blob.Alive
}

func (p *PlayerInfo) Die() {
	// p.mux.Lock()
	// defer p.mux.Unlock()
	p.Blob.Alive = false
}

func (p *PlayerInfo) GetMass() int32 {
	// p.mux.Lock()
	// defer p.mux.Unlock()
	return p.Blob.Mass
}

func (p *PlayerInfo) GetBlob() *Blob {
	// p.mux.Lock()
	// defer p.mux.Unlock()
	b := &Blob{
		Ip:    p.Blob.Ip,
		X:     p.Blob.X,
		Y:     p.Blob.Y,
		Mass:  p.Blob.Mass,
		Alive: p.Blob.Alive,
		Ver:   p.Blob.Ver,
		Seq:   p.Blob.Seq,
	}
	return b
}

func (p *PlayerInfo) UpdatePos(dx float64, dy float64) (float64, float64) {
	// p.mux.Lock()
	// defer p.mux.Unlock()
	updateBlob := &p.Blob
	updateBlob.Seq++

	updateBlob.X += dx
	updateBlob.Y += dy

	//constrain movement for now
	if updateBlob.X > float64(Conf.MAP_WIDTH-1) {
		updateBlob.X = float64(Conf.MAP_WIDTH-1)
	} else if updateBlob.X < 0 {
		updateBlob.X = 0
	}

	if updateBlob.Y > float64(Conf.MAP_HEIGHT-1) {
		updateBlob.Y = float64(Conf.MAP_HEIGHT-1)
	} else if updateBlob.Y < 0 {
		updateBlob.Y = 0
	}
	// log.Println(updateBlob.X, updateBlob.Y)
	return updateBlob.X, updateBlob.Y
}

func (p *PlayerInfo) IncrementMass(deltaMass int32) {
	// p.mux.Lock()
	// defer p.mux.Unlock()

	//Poison food
	prob := rand.Intn(Conf.POISON_PROB)
	if prob == 0 {
		deltaMass *= -1
	}

	p.Blob.Mass = max(p.Blob.Mass+deltaMass, Conf.STARTING_MASS)
}

func (p *PlayerInfo) UpdateSeq() {
	// p.mux.Lock()
	// defer p.mux.Unlock()

	p.Blob.Seq = time.Now().UnixNano()
}


// Returns a list of region ids
func (p *PlayerInfo) GetAOI() []uint32 {
	// p.mux.Lock()
	// defer p.mux.Unlock()
	diameter := GetRadiusFromMass(p.Blob.Mass) * 2.0
	zoom_factor := Conf.ZOOM/diameter + 0.3
	top_left_x := p.Blob.X - float64(Conf.SCREEN_WIDTH)/zoom_factor/2
	top_left_y := p.Blob.Y - float64(Conf.SCREEN_HEIGHT)/zoom_factor/2
	bot_right_x := p.Blob.X + float64(Conf.SCREEN_WIDTH)/zoom_factor/2
	bot_right_y := p.Blob.Y + float64(Conf.SCREEN_HEIGHT)/zoom_factor/2
	//figure out which of the regions are in AOI
	start_region_x := max(0, int32(math.Floor(top_left_x/float64(Conf.REGION_MAP_WIDTH))))
	start_region_y := max(0, int32(math.Floor(top_left_y/float64(Conf.REGION_MAP_HEIGHT))))
	end_region_x := min(int32(math.Ceil(bot_right_x/float64(Conf.REGION_MAP_WIDTH))), int32(Conf.NREGION_WIDTH-1))
	end_region_y := min(int32(math.Ceil(bot_right_y/float64(Conf.REGION_MAP_HEIGHT))), int32(Conf.NREGION_HEIGHT-1))

	regionIds := make([]uint32, 0)
	// log.Println("start & end x", start_region_x, end_region_x)
	// log.Println("start & end y", start_region_y, end_region_y)

	for x := start_region_x; x <= end_region_x; x++ {
		for y := start_region_y; y <= end_region_y; y++ {
			// log.Println("Region: ", int(x), int(y))
			regionIds = append(regionIds, getRegionID(uint16(x), uint16(y)))
		}
	}
	return regionIds
}

func getRegionID(x, y uint16) uint32 {
	return uint32(x)<<16 | uint32(y)
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
