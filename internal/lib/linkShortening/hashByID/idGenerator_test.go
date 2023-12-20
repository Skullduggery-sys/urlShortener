package hashByID

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestgetIDZero(t *testing.T) {
	gen := newIDGenerator(0)
	res := 100000
	wg := sync.WaitGroup{}
	wg.Add(res)
	for i := 0; i < res; i++ {
		go func() {
			gen.getID()
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, uint64(res), gen.id)
}

func TestgetIDNotZeros(t *testing.T) {
	var init uint64 = 2131
	gen := newIDGenerator(init)
	res := 100000
	wg := sync.WaitGroup{}
	wg.Add(res)
	for i := 0; i < res; i++ {
		go func() {
			gen.getID()
			wg.Done()
		}()
	}
	wg.Wait()
	assert.Equal(t, uint64(res)+init, gen.id)
}
