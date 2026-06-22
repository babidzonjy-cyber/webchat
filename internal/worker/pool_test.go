package worker

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestPool_Submit(t *testing.T) {
	pool := NewPool(3)
	defer pool.Shutdown()

	var counter int64
	var wg sync.WaitGroup

	for range 10 {
		wg.Add(1)
		pool.Submit(func() {
			defer wg.Done()
			atomic.AddInt64(&counter, 1)
		})
	}
	wg.Wait()

	if counter != 10 {
		t.Errorf("expected 10, got %d", counter)
	}
}

func TestPool_Concurrent(t *testing.T) {
	pool := NewPool(3)
	defer pool.Shutdown()

	var maxConcurrent int64
	var current int64

	for range 10 {
		pool.Submit(func() {
			cur := atomic.AddInt64(&current, 1)
			defer atomic.AddInt64(&current, -1)

			for {
				old := atomic.LoadInt64(&maxConcurrent)
				if cur <= old || atomic.CompareAndSwapInt64(&maxConcurrent, old, cur) {
					break
				}
			}
			time.Sleep(10 * time.Millisecond)
		})
	}

	time.Sleep(200 * time.Millisecond)
	if atomic.LoadInt64(&maxConcurrent) < 2 {
		t.Errorf("expected at least 2 concurrent, got %d", maxConcurrent)
	}
}

func TestPool_Shutdown(t *testing.T) {
	pool := NewPool(2)

	var done int64

	for range 10 {
		pool.Submit(func() {
			time.Sleep(10 * time.Millisecond)
			atomic.AddInt64(&done, 1)
		})
	}

	pool.Shutdown()

	if done != 10 {
		t.Errorf("expected 10, got %d", done)
	}
}
