package region

type RegionInfo struct {
	foodTree *quadtree.Quadtree
	playersBelongs  map[string]*player.Player
	playersPresents map[string]*player.Player
	mux      sync.Mutex
	x        uint16   
	y        uint16   
	hash     uint32   
}
