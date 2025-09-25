package sync

import (
	isync "internal/sync"
)

// NOTE: Not using generics for the sake of simplicity in the slides

/******************************************************************************/
/*                                  MyMutex1                                  */
/******************************************************************************/

type MyMutex1 struct {
	_  noCopy
	mu *isync.MyMutex1
}

func NewMyMutex1() *MyMutex1 {
	mu := isync.NewMyMutex1()
	return &MyMutex1{mu: mu}
}

func (m *MyMutex1) Lock() {
	m.mu.Lock()
}

func (m *MyMutex1) Unlock() {
	m.mu.Unlock()
}

/******************************************************************************/
/*                                  MyMutex2                                  */
/******************************************************************************/

type MyMutex2 struct {
	_  noCopy
	mu isync.MyMutex2
}

func (m *MyMutex2) TryLock() bool {
	return m.mu.TryLock()
}

func (m *MyMutex2) Lock() {
	m.mu.Lock()
}

func (m *MyMutex2) Unlock() {
	m.mu.Unlock()
}

/******************************************************************************/
/*                                  MyMutex3                                  */
/******************************************************************************/

type MyMutex3 struct {
	_  noCopy
	mu isync.MyMutex3
}

func (m *MyMutex3) TryLock() bool {
	return m.mu.TryLock()
}

func (m *MyMutex3) Lock() {
	m.mu.Lock()
}

func (m *MyMutex3) Unlock() {
	m.mu.Unlock()
}

/******************************************************************************/
/*                                  MyMutex4                                  */
/******************************************************************************/

type MyMutex4 struct {
	_  noCopy
	mu isync.MyMutex4
}

func (m *MyMutex4) TryLock() bool {
	return m.mu.TryLock()
}

func (m *MyMutex4) Lock() {
	m.mu.Lock()
}

func (m *MyMutex4) Unlock() {
	m.mu.Unlock()
}

/******************************************************************************/
/*                                  MyMutex5                                  */
/******************************************************************************/

type MyMutex5 struct {
	_  noCopy
	mu isync.MyMutex5
}

func (m *MyMutex5) TryLock() bool {
	return m.mu.TryLock()
}

func (m *MyMutex5) Lock() {
	m.mu.Lock()
}

func (m *MyMutex5) Unlock() {
	m.mu.Unlock()
}
