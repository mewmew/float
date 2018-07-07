package binary16

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFromBits(t *testing.T) {
	// from: https://en.wikipedia.org/wiki/Half-precision_floating-point_format#Half_precision_examples
	golden := []struct {
		bits uint16
		want float64
	}{
		// 0 01111 0000000000 = 1
		{bits: 0x3C00, want: 1},
		// 0 01111 0000000001 = 1 + 2^(-10) = 1.0009765625 (next smallest float after 1)
		{bits: 0x3C01, want: 1.0009765625},
		// 1 10000 0000000000 = -2
		{bits: 0xC000, want: -2},
		// 0 11110 1111111111 = 65504 (max half precision)
		{bits: 0x7BFF, want: 65504},
		// 0 00001 0000000000 = 2^(-14) ~= 6.10352 * 10^(-5) (minimum positive normal)
		{bits: 0x0400, want: math.Pow(2, -14)},
		// 0 00000 0000000001 = 2^(-24) ~= 5.96046 * 10^(-8) (minimum positive subnormal)
		{bits: 0x0001, want: math.Pow(2, -24)},
		// 0 00000 0000000000 = 0
		{bits: 0x0000, want: 0},
		// 1 00000 0000000000 = −0
		{bits: 0x8000, want: math.Copysign(0, -1)},
		// 0 11111 0000000000 = infinity
		{bits: 0x7C00, want: math.Inf(1)},
		// 1 11111 0000000000 = -infinity
		{bits: 0xFC00, want: math.Inf(-1)},
		// 0 01101 0101010101 = 0.333251953125 ~= 1/3
		{bits: 0x3555, want: 0.333251953125},
	}
	for _, g := range golden {
		f := NewFromBits(g.bits)
		got := f.Float64()
		wantBits := math.Float64bits(g.want)
		gotBits := math.Float64bits(got)
		//fmt.Printf("bits: 0x%04X (%v)\n", g.bits, g.want)
		if wantBits != gotBits {
			t.Errorf("0x%04X: number mismatch; expected 0x%08X (%v), got 0x%08X (%v)", g.bits, wantBits, g.want, gotBits, got)
		}
	}
}
