package sync

import (
	isync "internal/sync"
)

type MyMutex struct {
	_ noCopy

	// mu isync.MyMutex1
	mu isync.MyMutex2
}

// var _ Locker = &isync.MyMutex1{}
var _ Locker = &isync.MyMutex2{}

// func NewMyMutex() isync.MyMutex1 {
// 	return isync.NewMyMutex1()
// }

func (m *MyMutex) Lock() {
	m.mu.Lock()
}

func (m *MyMutex) Unlock() {
	m.mu.Unlock()
}
