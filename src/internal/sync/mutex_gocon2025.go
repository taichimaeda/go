package sync

import (
	// "internal/race"
	"sync/atomic"
	// "unsafe"
)

const (
	myMutexLocked      = 1 << iota
	myMutexWaiterShift = iota
)

/******************************************************************************/
/*                                  MyMutex1                                  */
/******************************************************************************/

type MyMutex1 struct {
	sema uint32
}

func NewMyMutex1() MyMutex1 {
	return MyMutex1{sema: 1} // need to init sema to 1
}

// NOTE: No TryLock() possible

func (m *MyMutex1) Lock() {
	println("Locking MyMutex1...") // using builtin println() to prevent cyclic deps
	defer println("Locking MyMutex1 complete!")

	queueLifo := false
	skipframes := 1 // skip 1 caller from stack trace (sync.MyMutex.Lock())
	runtime_SemacquireMutex(&m.sema, queueLifo, skipframes)
}

func (m *MyMutex1) Unlock() {
	println("Unlocking MyMutex1...")
	defer println("Unlocking MyMutex1 complete!")

	handoff := false
	skipframes := 1
	runtime_Semrelease(&m.sema, handoff, skipframes)
}

/******************************************************************************/
/*                                  MyMutex2                                  */
/******************************************************************************/

type MyMutex2 struct {
	state int32
	sema  uint32
}

func (m *MyMutex2) TryLock() bool {
	println("Trying to lock MyMutex2...")
	defer println("Trying to lock MyMutex2 complete!")

	if atomic.SwapInt32(&m.state, myMutexLocked) != 0 {
		return false
	}
	return true
}

func (m *MyMutex2) Lock() {
	println("Locking MyMutex2...")
	defer println("Locking MyMutex2 complete!")

	for atomic.SwapInt32(&m.state, myMutexLocked) != 0 {
		queueLifo := false
		skipframes := 1
		runtime_SemacquireMutex(&m.sema, queueLifo, skipframes)
	}
}

func (m *MyMutex2) Unlock() {
	println("Unlocking MyMutex2...")
	defer println("Unlocking MyMutex2 complete!")

	atomic.StoreInt32(&m.state, 0)
	handoff := false
	skipframes := 1
	runtime_Semrelease(&m.sema, handoff, skipframes)
}

/******************************************************************************/
/*                                  MyMutex3                                  */
/******************************************************************************/

type MyMutex3 struct {
	state int32
	sema  uint32
}

func (m *MyMutex3) TryLock() bool {
	println("Trying to lock MyMutex3...")
	defer println("Trying to lock MyMutex3 complete!")

	if atomic.SwapInt32(&m.state, myMutexLocked) != 0 {
		return false
	}
	return true
}

func (m *MyMutex3) Lock() {
	println("Locking MyMutex3...")
	defer println("Locking MyMutex3 complete!")

	iter := 0
	for atomic.SwapInt32(&m.state, myMutexLocked) != 0 {
		if runtime_canSpin(iter) {
			runtime_doSpin()
			iter++
			continue
		}
		queueLifo := false
		skipframes := 1
		runtime_SemacquireMutex(&m.sema, queueLifo, skipframes)
	}
}

func (m *MyMutex3) Unlock() {
	println("Unlocking MyMutex3...")
	defer println("Unlocking MyMutex3 complete!")

	atomic.StoreInt32(&m.state, 0)
	handoff := false
	skipframes := 1
	runtime_Semrelease(&m.sema, handoff, skipframes)
}

/******************************************************************************/
/*                                  MyMutex4                                  */
/******************************************************************************/

// TODO: This is the same impl as MyMutex3

type MyMutex4 struct {
	state int32
	sema  uint32
}

func (m *MyMutex4) TryLock() bool {
	println("Trying to lock MyMutex4...")
	defer println("Trying to lock MyMutex4 complete!")

	old := m.state
	if old&myMutexLocked != 0 {
		return false
	}
	if !atomic.CompareAndSwapInt32(&m.state, old, old|myMutexLocked) { // could be still unlocked but waiter count decremented by Unlock()
		return false
	}
	return true
}

func (m *MyMutex4) Lock() {
	println("Locking MyMutex4...")
	defer println("Locking MyMutex4 complete!")

	// fast path
	if atomic.CompareAndSwapInt32(&m.state, 0, myMutexLocked) {
		return
	}

	// slow path
	iter := 0
	old := m.state // not atomic but ok
	for {
		if old&myMutexLocked == myMutexLocked && runtime_canSpin(iter) {
			runtime_doSpin()
			iter++
			old = m.state
			continue
		}
		new := old
		new |= myMutexLocked
		if old&myMutexLocked == myMutexLocked {
			new += 1 << myMutexWaiterShift
		}
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			if old&myMutexLocked == 0 {
				break
			}
			queueLifo := false
			skipframes := 1
			runtime_SemacquireMutex(&m.sema, queueLifo, skipframes)
			iter = 0
		}
		old = m.state
	}
}

func (m *MyMutex4) Unlock() {
	println("Unlocking MyMutex4...")
	defer println("Unlocking MyMutex4 complete!")

	// fast path
	new := atomic.AddInt32(&m.state, -myMutexLocked)
	if new == 0 {
		return
	}
	if (new+mutexLocked)&mutexLocked == 0 { // add back myMutexLocked in case it was not set initially
		fatal("Unlocked unlocked MyMutex4!")
	}

	// slow path
	old := new
	for {
		if old>>myMutexWaiterShift == 0 || old&myMutexLocked != 0 {
			return
		}
		new = old - 1<<myMutexWaiterShift
		if atomic.CompareAndSwapInt32(&m.state, old, new) {
			handoff := false
			skipframes := 1
			runtime_Semrelease(&m.sema, handoff, skipframes)
		}
		old = m.state
	}
}
