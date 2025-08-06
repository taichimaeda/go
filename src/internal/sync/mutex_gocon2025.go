package sync

import (
	// "internal/race"
	"sync/atomic"
	// "unsafe"
)

const (
	myMutexLocked = 1 << iota
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
	queueLifo := false
	skipframes := 1 // skip 1 caller from stack trace (sync.MyMutex.Lock())
	runtime_SemacquireMutex(&m.sema, queueLifo, skipframes)
	println("Locking MyMutex1 complete!")
}

func (m *MyMutex1) Unlock() {
	println("Unlocking MyMutex1...")
	handoff := false
	skipframes := 1
	runtime_Semrelease(&m.sema, handoff, skipframes)
	println("Unlocking MyMutex1 complete!")
}

/******************************************************************************/
/*                                  MyMutex2                                  */
/******************************************************************************/

type MyMutex2 struct {
	state int32
	sema  uint32
}

func (m *MyMutex2) TryLock() bool {
	if atomic.SwapInt32(&m.state, mutexLocked) != 0 {
		return false
	}
	return true
}

func (m *MyMutex2) Lock() {
	println("Locking MyMutex2...")
	for atomic.SwapInt32(&m.state, myMutexLocked) != 0 {
		queueLifo := false
		skipframes := 1
		runtime_SemacquireMutex(&m.sema, queueLifo, skipframes)
	}
	println("Locking MyMutex2 complete!")
}

func (m *MyMutex2) Unlock() {
	println("Unlocking MyMutex2...")
	atomic.StoreInt32(&m.state, 0)
	handoff := false
	skipframes := 1
	runtime_Semrelease(&m.sema, handoff, skipframes)
	println("Unlocking MyMutex2 complete!")
}

/******************************************************************************/
/*                                  MyMutex3                                  */
/******************************************************************************/

type MyMutex3 struct {
	state int32
	sema  uint32
}

func (m *MyMutex3) TryLock() bool {
	if atomic.SwapInt32(&m.state, mutexLocked) != 0 {
		return false
	}
	return true
}

func (m *MyMutex3) Lock() {
	println("Locking MyMutex3...")
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
	println("Locking MyMutex3 complete!")
}

func (m *MyMutex3) Unlock() {
	println("Unlocking MyMutex3...")
	atomic.StoreInt32(&m.state, 0)
	handoff := false
	skipframes := 1
	runtime_Semrelease(&m.sema, handoff, skipframes)
	println("Unlocking MyMutex3 complete!")
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
	if atomic.SwapInt32(&m.state, mutexLocked) != 0 {
		return false
	}
	return true
}

func (m *MyMutex4) Lock() {
	println("Locking MyMutex4...")
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
	println("Locking MyMutex4 complete!")
}

func (m *MyMutex4) Unlock() {
	println("Unlocking MyMutex4...")
	atomic.StoreInt32(&m.state, 0)
	handoff := false
	skipframes := 1
	runtime_Semrelease(&m.sema, handoff, skipframes)
	println("Unlocking MyMutex4 complete!")
}
