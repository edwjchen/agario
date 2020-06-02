package player

import (
	"sync"
	"io/ioutil"
	"strconv"
	"time"
	. "peer_to_peer/common"
	// "log"
)

type PlayerInfo struct {
	Blob       Blob
	LastUpdate time.Time
	mux        sync.Mutex
	addr       string
}

func (p *PlayerInfo) InitIP(addr string) {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.addr = addr
}

func (p *PlayerInfo) NewBlob() (string, float64, float64, int32, int32) {
	// TODO change to map
	p.mux.Lock()
	defer p.mux.Unlock()

	var x float64 = 400.0
	var y float64 = 200.0

	var ver int
	verBytes, err := ioutil.ReadFile(VER_FILE)
	if err != nil {
		ver = 0
	} else {
		ver, _ = strconv.Atoi(string(verBytes))
		ver++
	}

	writeVer := []byte(strconv.Itoa(ver))
	ioutil.WriteFile(VER_FILE, writeVer, 0666) 
	
	p.Blob = Blob{Ip: p.addr, X:x, Y: y, Mass: STARTING_MASS, Alive: true, Ver: int32(ver)}
	return p.addr, x, y, STARTING_MASS, int32(ver)
}

func (p *PlayerInfo) GetAlive() bool {
	p.mux.Lock()
	defer p.mux.Unlock()
	return p.Blob.Alive
}

func (p *PlayerInfo) GetMass() int32 {
	p.mux.Lock()
	defer p.mux.Unlock()
	return p.Blob.Mass
}

func (p *PlayerInfo) GetBlob() *Blob {
	p.mux.Lock()
	defer p.mux.Unlock()
	b := &Blob{
		Ip: p.Blob.Ip,
		X: p.Blob.X,
		Y: p.Blob.Y,
		Mass: p.Blob.Mass,
		Alive: p.Blob.Alive,
		Ver: p.Blob.Ver,
	}
	return b
}

func (p *PlayerInfo) UpdatePos(dx float64, dy float64) (float64, float64) {
	p.mux.Lock()
	defer p.mux.Unlock()
	updateBlob := &p.Blob

	updateBlob.X += dx
	updateBlob.Y += dy

	//constrain movement for now
	if updateBlob.X > MAP_WIDTH {
		updateBlob.X = MAP_WIDTH
	} else if updateBlob.X < 0 {
		updateBlob.X = 0
	}

	if updateBlob.Y > MAP_HEIGHT {
		updateBlob.Y = MAP_HEIGHT
	} else if updateBlob.Y < 0 {
		updateBlob.Y = 0
	}
	// log.Println(updateBlob.X, updateBlob.Y)
	return updateBlob.X, updateBlob.Y
}

func (p *PlayerInfo) IncrementMass(deltaMass int32) {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.Blob.Mass += deltaMass
}

// Returns a list of region ids
func (p *PlayerInfo) GetAOI() []uint32 {
	p.mux.Lock()
	defer p.mux.Unlock()
	zoom_factor := float64(ZOOM/p.Blob.Mass)+0.3
	top_left_x := p.Blob.X - float64(SCREEN_WIDTH)/zoom_factor/2
	top_left_y := p.Blob.Y - float64(SCREEN_HEIGHT)/zoom_factor/2
	bot_right_x := p.Blob.X + float64(SCREEN_WIDTH)/zoom_factor/2
	bot_right_y := p.Blob.Y + float64(SCREEN_HEIGHT)/zoom_factor/2
	//figure out which of the regions are in AOI
	start_region_x := int(top_left_x/REGION_MAP_WIDTH)
	start_region_y := int(top_left_y/REGION_MAP_HEIGHT)
	end_region_x := int(bot_right_x/REGION_MAP_WIDTH)
	end_region_y := int(bot_right_y/REGION_MAP_HEIGHT)

	regionIds := make([]uint32, 0)
	// log.Println("start & end x", start_region_x, end_region_x)
	// log.Println("start & end y", start_region_y, end_region_y)

	for x := start_region_x; x <= end_region_x; x++ {
		for y := start_region_y; y <= end_region_y; y++ {
			// log.Println("Region: ", int(x), int(y))
			regionIds = append(regionIds, getRegionID(uint16(x),uint16(y)))
		} 
	}
	return regionIds
}

func getRegionID(x, y uint16) uint32 {
	return uint32(x) << 16 | uint32(y)
}