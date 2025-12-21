package test

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/foxie-io/ng"
	"github.com/stretchr/testify/assert"
)

type Key struct {
	id any
}

func (k Key) PayloadKey() string {
	return fmt.Sprintf("Key-%v", k.id)
}

func TestRequestContextPoolRacing(t *testing.T) {
	assert := assert.New(t)

	const goroutines = 100
	const iterations = 1000

	var wg sync.WaitGroup

	checker := sync.Map{}

	counter := atomic.Int64{}
	newId := func() int64 {
		return counter.Add(1)
	}

	// Concurrently acquire, use, and release contexts with potential racing
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for j := 0; j < iterations; j++ {
				_, ctx := ng.AcquireContext(context.Background())
				defer ctx.Clear()

				id := newId()

				ctx.Store(Key{"id"}, id)
				val, ok := ctx.Load(Key{"id"})

				assert.True(ok)
				assert.Equal(id, val, " goroutine %d, iteration %d", i, j)

				// Atomically check and store context ID
				if _, loaded := checker.LoadOrStore(val, struct{}{}); loaded {
					assert.Fail("Duplicate context ID detected", id)
				}
			}
		}(i)
	}

	wg.Wait()
}

// BenchmarkAcquireContext-8   	18629553	        66.67 ns/op	     208 B/op	       2 allocs/op
func BenchmarkAcquireContext(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ng.AcquireContext(context.Background())
	}
}
