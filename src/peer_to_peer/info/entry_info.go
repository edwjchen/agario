package info

import (
	"math/rand"
	"sync"
	"time"
)
type EntryInfo struct {
	mux        sync.Mutex
	MinPlayers int
	MaxPlayers int 
	CurrNodes  []string
}

func (e *EntryInfo) AddNode(node string) {
	e.mux.Lock()
	defer e.mux.Unlock()

	e.CurrNodes = append(e.CurrNodes, node)
}

func (e *EntryInfo) CanStart() bool {
	e.mux.Lock()
	defer e.mux.Unlock()

	return len(e.CurrNodes) >= e.MinPlayers
}

func (e *EntryInfo) GetNodeIP() string {
	e.mux.Lock()
	defer e.mux.Unlock()

	rand.Seed(time.Now().UnixNano())
	return e.CurrNodes[rand.Intn(len(e.CurrNodes))]
}

func (e *EntryInfo) ShouldCreate() bool {
	e.mux.Lock()
	defer e.mux.Unlock()
	
	return len(e.CurrNodes) < e.MinPlayers
}