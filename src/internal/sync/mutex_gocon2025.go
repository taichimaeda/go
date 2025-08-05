package sync

import (
	// "internal/race"
	"sync/atomic"
	// "unsafe"
)

const (
	myMutexLocked = 1 << iota
)

type MyMutex1 struct {
	sema uint32
}

func NewMyMutex1() MyMutex1 {
	return MyMutex1{sema: 1} // need to init sema to 1
}

func (m *MyMutex1) Lock() {
	println("Locking MyMutex...") // using builtin println() to prevent cyclic deps
	queueLifo := false
	skipframes := 1 // skip 1 caller from stack trace (sync.MyMutex.Lock())
	runtime_SemacquireMutex(&m.sema, queueLifo, skipframes)
	println("Locking MyMutex complete!")
}

func (m *MyMutex1) Unlock() {
	println("Unlocking MyMutex...")
	handoff := false
	skipframes := 1
	runtime_Semrelease(&m.sema, handoff, skipframes)
	println("Unlocking MyMutex complete!")
}

type MyMutex2 struct {
	state int32
	sema  uint32
}

// TODO: Add TryLock()

func (m *MyMutex2) Lock() {
	println("Locking MyMutex...")
	for atomic.SwapInt32(&m.state, myMutexLocked) != 0 {
		queueLifo := false
		skipframes := 1
		runtime_SemacquireMutex(&m.sema, queueLifo, skipframes)
	}
	println("Locking MyMutex complete!")
}

func (m *MyMutex2) Unlock() {
	println("Unlocking MyMutex...")
	atomic.StoreInt32(&m.state, 0)
	handoff := false
	skipframes := 1
	runtime_Semrelease(&m.sema, handoff, skipframes)
	println("Unlocking MyMutex complete!")
}

type MyMutex3 struct {
	state int32
	sema  uint32
}

// TODO: Add TryLock()

func (m *MyMutex3) Lock() {
	println("Locking MyMutex...")
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
	println("Locking MyMutex complete!")
}

func (m *MyMutex3) Unlock() {
	println("Unlocking MyMutex...")
	atomic.StoreInt32(&m.state, 0)
	handoff := false
	skipframes := 1
	runtime_Semrelease(&m.sema, handoff, skipframes)
	println("Unlocking MyMutex complete!")
}
