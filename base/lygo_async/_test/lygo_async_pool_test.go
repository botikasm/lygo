package _test

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_async"
	"sync"
	"testing"
)

func TestExample(t *testing.T) {
	pool := lygo_async.NewConcurrentPool(10)
	for i := 0; i < 1000; i++ {
		pool.Run(func() {
			// do some work
			fmt.Println(pool.Count())
		})
	}
	pool.Join()
}

func TestLimit(t *testing.T) {
	LIMIT := 10
	N := 100

	pool := lygo_async.NewConcurrentPool(LIMIT)
	m := map[int]bool{}
	lock := &sync.Mutex{}

	max := int32(0)
	for i := 0; i < N; i++ {
		x := i
		ticket := pool.Run(func() {
			lock.Lock()
			m[x] = true
			currentMax := pool.Count()
			if currentMax >= max {
				max = currentMax
			}
			lock.Unlock()
		})
		if ticket > LIMIT-1 {
			t.Errorf("expected max: %d, got %d", LIMIT, ticket)
		}
	}

	pool.Join()

	t.Log("results:", len(m))
	t.Log("max:", max)

	if len(m) != N {
		t.Error("invalid num of results", len(m))
	}

	if max > int32(LIMIT) || max == 0 {
		t.Error("invalid max", max)
	}
}
