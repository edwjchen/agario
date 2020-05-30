package player

import (
	"sync"
)

type PlayerInfo struct {
	blob Blob
	mux  sync.Mutex
	addr string
}

const STARTING_MASS = 20
const MAP_WIDTH = 10000
const MAP_HEIGHT = 10000

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
	p.blob = Blob{Name: p.addr, X: x, Y: y, Mass: STARTING_MASS}
	return p.addr, x, y, STARTING_MASS
}

func (p *PlayerInfo) GetAlive() bool {
	p.mux.Lock()
	defer p.mux.Lock()
	return p.blob.Alive
}

func (p *PlayerInfo) GetMass() int32 {
	p.mux.Lock()
	defer p.mux.Lock()
	return p.blob.Mass
}

func (p *PlayerInfo) UpdatePos(dx float64, dy float64) (float64, float64) {
	p.mux.Lock()
	defer p.mux.Unlock()
	updateBlob := p.blob

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
	// log.Println(name, " is at pos ", updateBlob.x, updateBlob.y)
	return updateBlob.X, updateBlob.Y
}

