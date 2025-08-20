// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strconv_test

import (
	"math"
	"math/rand"
	. "strconv"
	"testing"
	"time"
)

func randomFloat32FullRange() float32 {
	bits := rand.Uint32() // random 32-bit pattern
	return math.Float32frombits(bits)
}

func randomFloat64FullRange() float64 {
	bits := rand.Uint64() // random 64-bit pattern
	return math.Float64frombits(bits)
}

func TestFtoaDragonbox(t *testing.T) {
	const iter = 100000

	for i := 0; i < iter; i++ {
		var bitSize int
		var val float64
		switch rand.Intn(2) {
		case 0:
			bitSize = 32
			val = float64(randomFloat32FullRange())
		case 1:
			bitSize = 64
			val = randomFloat64FullRange()
		}

		output1, _ := RunDragonboxFtoa(val, bitSize)
		output2, _ := RunRyuFtoaShortest(val, bitSize)

		if output1 != output2 {
			t.Errorf("Mismatch:\nInput: %f\nDragonbox output: %s\nRyu output: %s", val, output1, output2)
		}
	}
}

func BenchmarkDragonboxFtoa(b *testing.B) {
	const numTests = 100000

	var total1, total2 time.Duration

	for b.Loop() {
		for i := 0; i < numTests; i++ {
			var bitSize int
			var val float64
			switch rand.Intn(2) {
			case 0:
				bitSize = 32
				val = float64(randomFloat32FullRange())
			case 1:
				bitSize = 64
				val = randomFloat64FullRange()
			}

			_, elapsed1 := RunDragonboxFtoa(val, bitSize)
			_, elapsed2 := RunRyuFtoaShortest(val, bitSize)
			total1 += elapsed1
			total2 += elapsed2
		}
	}

	b.Logf("Duration:\nDragonbox: %d\nRyu: %d", total1, total2)
}
