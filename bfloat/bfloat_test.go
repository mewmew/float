package bfloat

import (
	"math"
	"testing"
)

func TestNewFromBits(t *testing.T) {
	golden := []struct {
		bits uint16
		want float64
	}{
		// Special numbers.
		// 0 00000000 0000000 = 0
		{bits: 0, want: 0},
		// 1 00000000 0000000 = -0
		{bits: 0x8000, want: 1. / math.Inf(-1)},
		// 0 11111111 0000000 = +Inf
		{bits: 0x7f80, want: math.Inf(1)},
		// 1 11111111 0000000 = -Inf
		{bits: 0xff80, want: math.Inf(-1)},

		// FIXME: the following comment case correctly get into nan branch but calculation of big.Float is incorrect
		// // 0 11111111 0000001 = +NaN
		// {bits: 0x7f81, want: math.NaN()},
		// // 1 11111111 0000001 = -NaN
		// {bits: 0xff81, want: -math.NaN()},

		// from: https://en.wikipedia.org/wiki/Bfloat16_floating-point_format#Examples
		{bits: 0x3f80, want: 1},
		{bits: 0xc000, want: -2},
		{bits: 0x4049, want: 3.140625},
		{bits: 0x3eab, want: 0.333984375},

		// // 2^i
		{bits: 0x0001, want: math.Pow(2, -24)}, // 2^(-24)
		{bits: 0x0002, want: math.Pow(2, -23)}, // 2^(-23)
		{bits: 0x0004, want: math.Pow(2, -22)}, // 2^(-22)
		{bits: 0x0008, want: math.Pow(2, -21)}, // 2^(-21)
		{bits: 0x0010, want: math.Pow(2, -20)}, // 2^(-20)
		{bits: 0x0020, want: math.Pow(2, -19)}, // 2^(-19)
		{bits: 0x0040, want: math.Pow(2, -18)}, // 2^(-18)
		{bits: 0x0080, want: math.Pow(2, -17)}, // 2^(-17)
		{bits: 0x0100, want: math.Pow(2, -16)}, // 2^(-16)
		{bits: 0x0200, want: math.Pow(2, -15)}, // 2^(-15)
		{bits: 0x0400, want: math.Pow(2, -14)}, // 2^(-14)
		{bits: 0x0800, want: math.Pow(2, -13)}, // 2^(-13)
		{bits: 0x0C00, want: math.Pow(2, -12)}, // 2^(-12)
		{bits: 0x1000, want: math.Pow(2, -11)}, // 2^(-11)
		{bits: 0x1400, want: math.Pow(2, -10)}, // 2^(-10)
		{bits: 0x1800, want: math.Pow(2, -9)},  // 2^(-9)
		{bits: 0x1C00, want: math.Pow(2, -8)},  // 2^(-8)
		{bits: 0x2000, want: math.Pow(2, -7)},  // 2^(-7)
		{bits: 0x2400, want: math.Pow(2, -6)},  // 2^(-6)
		{bits: 0x2800, want: math.Pow(2, -5)},  // 2^(-5)
		{bits: 0x2C00, want: math.Pow(2, -4)},  // 2^(-4)
		{bits: 0x3000, want: math.Pow(2, -3)},  // 2^(-3)
		{bits: 0x3400, want: math.Pow(2, -2)},  // 2^(-2)
		{bits: 0x3800, want: math.Pow(2, -1)},  // 2^(-1)
		{bits: 0x3C00, want: math.Pow(2, 0)},   // 2^0
		{bits: 0x4000, want: math.Pow(2, 1)},   // 2^1
		{bits: 0x4400, want: math.Pow(2, 2)},   // 2^2
		{bits: 0x4800, want: math.Pow(2, 3)},   // 2^3
		{bits: 0x4C00, want: math.Pow(2, 4)},   // 2^4
		{bits: 0x5000, want: math.Pow(2, 5)},   // 2^5
		{bits: 0x5400, want: math.Pow(2, 6)},   // 2^6
		{bits: 0x5800, want: math.Pow(2, 7)},   // 2^7
		{bits: 0x5C00, want: math.Pow(2, 8)},   // 2^8
		{bits: 0x6000, want: math.Pow(2, 9)},   // 2^9
		{bits: 0x6400, want: math.Pow(2, 10)},  // 2^10
		{bits: 0x6800, want: math.Pow(2, 11)},  // 2^11
		{bits: 0x6C00, want: math.Pow(2, 12)},  // 2^12
		{bits: 0x7000, want: math.Pow(2, 13)},  // 2^13
		{bits: 0x7400, want: math.Pow(2, 14)},  // 2^14
		{bits: 0x7800, want: math.Pow(2, 15)},  // 2^15
	}
	for _, g := range golden {
		f := NewFromBits(g.bits)
		b, _ := f.Big()
		got, _ := b.Float64()
		wantBits := math.Float64bits(g.want)
		gotBits := math.Float64bits(got)
		if wantBits != gotBits {
			t.Errorf("0x%04X: number mismatch; expected 0x%016X (%v), got 0x%016X (%v)", g.bits, wantBits, g.want, gotBits, got)
		}
	}
}
