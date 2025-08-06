package sync

import (
	isync "internal/sync"
)

type MyMutex struct {
	_ noCopy

	// mu isync.MyMutex1
	// mu isync.MyMutex2
	// mu isync.MyMutex3
	// mu isync.MyMutex4
	// mu isync.MyMutex5
	mu isync.MyMutex6
}

var _ Locker = &MyMutex{}

// func NewMyMutex() isync.MyMutex1 {
// 	return isync.NewMyMutex1()
// }

func (m *MyMutex) TryLock() bool {
	return m.mu.TryLock()
}

func (m *MyMutex) Lock() {
	m.mu.Lock()
}

func (m *MyMutex) Unlock() {
	m.mu.Unlock()
}
