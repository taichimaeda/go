package sync_test

import (
	. "sync"
	"testing"
)

func TestMyMutex(t *testing.T) {
	var mu MyMutex
	mu.Lock()
	mu.Unlock()
}
