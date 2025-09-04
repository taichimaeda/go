package sync_test

import (
	. "sync"
	"testing"
	"time"
)

// NOTE: Not using generics for the sake of simplicity in the slides

/******************************************************************************/
/*                                  MyMutex1                                  */
/******************************************************************************/

func hammerMyMutex1(m *MyMutex1, loops int, done chan bool) {
	for i := 0; i < loops; i++ {
		m.Lock()
		m.Unlock()
	}
	done <- true
}

func TestMyMutex1(t *testing.T) {
	m := NewMyMutex1()

	m.Lock()
	m.Unlock()

	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go hammerMyMutex1(m, 1000, done)
	}
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatalf("can't acquire Mutex in 10 seconds")
		}
	}
}

func TestMyMutexFairness1(t *testing.T) {
	mu := NewMyMutex1()

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

/******************************************************************************/
/*                                  MyMutex2                                  */
/******************************************************************************/

func hammerMyMutex2(m *MyMutex2, loops int, done chan bool) {
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

func TestMyMutex2(t *testing.T) {
	var m MyMutex2

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
		go hammerMyMutex2(&m, 1000, done)
	}
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatalf("can't acquire Mutex in 10 seconds")
		}
	}
}

func TestMyMutexFairness2(t *testing.T) {
	var mu MyMutex2

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

/******************************************************************************/
/*                                  MyMutex3                                  */
/******************************************************************************/

func hammerMyMutex3(m *MyMutex3, loops int, done chan bool) {
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

func TestMyMutex3(t *testing.T) {
	var m MyMutex3

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
		go hammerMyMutex3(&m, 1000, done)
	}
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatalf("can't acquire Mutex in 10 seconds")
		}
	}
}

func TestMyMutexFairness3(t *testing.T) {
	var mu MyMutex3

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

/******************************************************************************/
/*                                  MyMutex4                                  */
/******************************************************************************/

func hammerMyMutex4(m *MyMutex4, loops int, done chan bool) {
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

func TestMyMutex4(t *testing.T) {
	var m MyMutex4

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
		go hammerMyMutex4(&m, 1000, done)
	}
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatalf("can't acquire Mutex in 10 seconds")
		}
	}
}

func TestMyMutexFairness4(t *testing.T) {
	var mu MyMutex4

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

/******************************************************************************/
/*                                  MyMutex5                                  */
/******************************************************************************/

func hammerMyMutex5(m *MyMutex5, loops int, done chan bool) {
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

func TestMyMutex5(t *testing.T) {
	var m MyMutex5

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
		go hammerMyMutex5(&m, 1000, done)
	}
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatalf("can't acquire Mutex in 10 seconds")
		}
	}
}

func TestMyMutexFairness5(t *testing.T) {
	var mu MyMutex5

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

/******************************************************************************/
/*                                  MyMutex6                                  */
/******************************************************************************/

func hammerMyMutex6(m *MyMutex6, loops int, done chan bool) {
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

func TestMyMutex6(t *testing.T) {
	var m MyMutex6

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
		go hammerMyMutex6(&m, 1000, done)
	}
	for i := 0; i < 10; i++ {
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			t.Fatalf("can't acquire Mutex in 10 seconds")
		}
	}
}

func TestMyMutexFairness6(t *testing.T) {
	var mu MyMutex6

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
