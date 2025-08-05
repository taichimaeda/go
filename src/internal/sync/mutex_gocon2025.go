package sync

// import (
// 	"internal/race"
// 	"sync/atomic"
// 	"unsafe"
// )

type MyMutex struct {
	state int32
	sema  uint32
}

func NewMyMutex() MyMutex {
	return MyMutex{
		sema: 1,
	}
}

func (m *MyMutex) Lock() {
	println("Locking MyMutex...") // using builtin println() to prevent cyclic deps
	queueLifo := false
	skipframes := 1 // skip 1 caller from stack trace (sync.MyMutex.Lock())
	runtime_SemacquireMutex(&m.sema, queueLifo, skipframes)
	println("Locking MyMutex complete!")
}

func (m *MyMutex) Unlock() {
	println("Unlocking MyMutex...")
	handoff := false
	skipframes := 1
	runtime_Semrelease(&m.sema, handoff, skipframes)
	println("Unlocking MyMutex complete!")
}
