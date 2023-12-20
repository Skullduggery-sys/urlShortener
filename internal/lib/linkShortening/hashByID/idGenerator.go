package hashByID

import "sync"

type idGenerator struct {
	id uint64
	mu sync.RWMutex
}

func newIDGenerator(id uint64) *idGenerator {
	return &idGenerator{
		id: id,
		mu: sync.RWMutex{},
	}
}

func (gen *idGenerator) getID() uint64 {
	gen.mu.Lock()
	defer gen.mu.Unlock()
	prev := gen.id
	gen.id++
	return prev
}
