// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strconv_test

import (
	"math"
	"math/rand"
	. "strconv"
	"testing"
)

func randomFloat32FullRange() float32 {
	bits := rand.Uint32() // random 64-bit pattern
	return math.Float32frombits(bits)
}

func randomFloat64FullRange() float64 {
	bits := rand.Uint64() // random 64-bit pattern
	return math.Float64frombits(bits)
}

func TestFtoaDragonbox(t *testing.T) {
	for i := 0; i < 100000; i++ {
		var bitSize int
		if rand.Intn(2) == 0 {
			bitSize = 32
		} else {
			bitSize = 64
		}
		val32 := randomFloat32FullRange()
		val64 := randomFloat64FullRange()

		if err := CompareDragonboxFtoaAndRyuShortestFtoa(bitSize, val32, val64); err != nil {
			t.Errorf("Mismatch:\nInput: %f\nRyu output: %s\nDragonbox output: %s",
				err.Val, err.RyuOutput, err.DragonboxOutput)
		}
	}
}

func BenchmarkDragonboxFtoa(b *testing.B) {
	for b.Loop() {
		var bitSize int
		if rand.Intn(2) == 0 {
			bitSize = 32
		} else {
			bitSize = 64
		}
		val32 := randomFloat32FullRange()
		val64 := randomFloat64FullRange()

		RunDragonboxFtoa(bitSize, val32, val64)
	}
}

func BenchmarkRyuFtoaShortest(b *testing.B) {
	for b.Loop() {
		var bitSize int
		if rand.Intn(2) == 0 {
			bitSize = 32
		} else {
			bitSize = 64
		}
		val32 := randomFloat32FullRange()
		val64 := randomFloat64FullRange()

		RunRyuFtoaShortest(bitSize, val32, val64)
	}
}
