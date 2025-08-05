package sync_test

import (
	. "sync"
	"testing"
)

func TestMyMutex(t *testing.T) {
	mu := NewMyMutex()
	mu.Lock()
	mu.Unlock()
}
