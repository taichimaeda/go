package sync_test

import (
	. "sync"
	"testing"
	"time"
)

func TestMyMutex(t *testing.T) {
	mu := NewMyMutex()
	mu.Lock()
	mu.Unlock()
}

// TODO: Add equivalent to TestMutex
// TODO: Add equivalent to TestMutexMisuse

func TestMyMutexFairness(t *testing.T) {
	mu := NewMyMutex()
	stop := make(chan bool)
	defer close(stop)
	go func() {
		for {
			mu.Lock()
			time.Sleep(100 * time.Microsecond)
			mu.Unlock()
			select {
			case <-stop:
				return
			default:
			}
		}
	}()
	done := make(chan bool, 1)
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(100 * time.Microsecond)
			mu.Lock()
			mu.Unlock()
		}
		done <- true
	}()
	select {
	case <-done:
	case <-time.After(10 * time.Second):
		t.Fatalf("can't acquire Mutex in 10 seconds")
	}
}
