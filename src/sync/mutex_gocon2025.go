package sync

import (
	isync "internal/sync"
)

type MyMutex struct {
	_ noCopy

	mu isync.MyMutex
}

type MyLocker interface {
	Lock()
	Unlock()
}

func NewMyMutex() MyMutex {
	return MyMutex{
		mu: isync.NewMyMutex(),
	}
}

func (m *MyMutex) Lock() {
	m.mu.Lock()
}

func (m *MyMutex) Unlock() {
	m.mu.Unlock()
}
