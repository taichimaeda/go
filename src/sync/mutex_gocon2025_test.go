package sync_test

import (
	"runtime"
	. "sync"
	"testing"
	"time"
)

func HammerMyMutex(m *MyMutex, loops int, done chan bool) {
	for i := 0; i < loops; i++ {
		if i%3 == 0 {
			if m.TryLock() {
				m.Unlock()
			}
			continue
		}
		m.Lock()
		m.Unlock()
	}
	done <- true
}

func TestMyMutex(t *testing.T) {
	if n := runtime.SetMutexProfileFraction(1); n != 0 {
		t.Logf("got mutexrate %d expected 0", n)
	}
	defer runtime.SetMutexProfileFraction(0)

	// m := NewMyMutex()
	var m MyMutex

	m.Lock()
	if m.TryLock() {
		t.Fatalf("TryLock succeeded with mutex locked")
	}
	m.Unlock()
	if !m.TryLock() {
		t.Fatalf("TryLock failed with mutex unlocked")
	}
	m.Unlock()

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go HammerMyMutex(&m, 1000, done)
	}
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatalf("can't acquire Mutex in 10 seconds")
		}
	}
}

func TestMyMutexFairness(t *testing.T) {
	// mu := NewMyMutex()
	var mu MyMutex
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
