package player

import (
	"sync"
	. "peer_to_peer/server/player_pb"
	// "log"
)

type PlayerInfo struct {
	blob Blob
	mux  sync.Mutex
	addr string
}

func (p *PlayerInfo) InitIP(addr string) {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.addr = addr
}

func (p *PlayerInfo) NewBlob() (string, float64, float64, int32) {
	// TODO change to map
	p.mux.Lock()
	defer p.mux.Unlock()

	var x float64 = 400.0
	var y float64 = 200.0
	p.blob = Blob{Name: p.addr, X:x, Y: y, Mass: STARTING_MASS, Alive: true}
	return p.addr, x, y, STARTING_MASS
}

func (p *PlayerInfo) GetAlive() bool {
	p.mux.Lock()
	defer p.mux.Unlock()
	return p.blob.Alive
}

func (p *PlayerInfo) GetMass() int32 {
	p.mux.Lock()
	defer p.mux.Unlock()
	return p.blob.Mass
}

func (p *PlayerInfo) GetBlob() *Blob {
	p.mux.Lock()
	defer p.mux.Unlock()
	b := &Blob{
		Name: p.blob.Name,
		X: p.blob.X,
		Y: p.blob.Y,
		Mass: p.blob.Mass,
		Alive: p.blob.Alive,
	}
	return b
}

func (p *PlayerInfo) UpdatePos(dx float64, dy float64) (float64, float64) {
	p.mux.Lock()
	defer p.mux.Unlock()
	updateBlob := &p.blob

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

	p.blob.Mass += deltaMass
}

// Returns a list of regions
func (p *PlayerInfo) GetAOI() {

}
