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

func (m *MyMutex) Lock() {
	// using builtin println() to prevent cyclic deps
	println("Locking MyMutex")
}

func (m *MyMutex) Unlock() {
	println("Unlocking MyMutex")
}
