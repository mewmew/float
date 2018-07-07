package binary16

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewFromBits(t *testing.T) {
	golden := []struct {
		bits uint16
		want float64
	}{
		// Special numbers.

		// +NaN
		{bits: 0x7E00, want: math.NaN()},
		// -NaN
		{bits: 0xFE00, want: -math.NaN()},
		// +Inf
		{bits: 0x7C00, want: math.Inf(1)},
		// -Inf
		{bits: 0xFC00, want: math.Inf(-1)},
		// +0
		{bits: 0x0000, want: 0.0},
		// -0
		{bits: 0x8000, want: math.Copysign(0.0, -1)},

		// From https://reviews.llvm.org/rL237161

		// Normalized numbers.
		{bits: 0x3800, want: 0.5},
		{bits: 0xB800, want: -0.5},
		{bits: 0x3E00, want: 1.5},
		{bits: 0xBE00, want: -1.5},
		{bits: 0x4100, want: 2.5},
		{bits: 0xC100, want: -2.5},

		// Denormalized numbers.
		{bits: 0x0010, want: float64FromString("0x1.0p-20")},
		{bits: 0x0001, want: float64FromString("0x1.0p-24")},
		{bits: 0x8001, want: float64FromString("-0x1.0p-24")},
		//{bits: 0x0001, want: float64FromString("0x1.5p-25")},

		// Rounding.
		// TODO: Handle rounding.
		//{bits: 0x4248, want: 3.14},
		//{bits: 0xC248, want: -3.14},
		//{bits: 0x4248, want: 3.1415926535},
		//{bits: 0xC248, want: -3.1415926535},
		//{bits: 0x7C00, want: float64FromString("0x1.987124876876324p+100")},
		{bits: 0x6E62, want: float64FromString("0x1.988p+12")},
		{bits: 0x3C00, want: float64FromString("0x1.0p+0")},
		{bits: 0x0400, want: float64FromString("0x1.0p-14")},
		// rounded to zero
		//{bits: 0x0000, want: float64FromString("0x1.0p-25")},
		//{bits: 0x8000, want: float64FromString("-0x1.0p-25")},
		// max (precise)
		{bits: 0x7BFF, want: 65504.0},
	}

	for _, g := range golden {
		f := NewFromBits(g.bits)
		got := f.Float64()
		wantBits := math.Float64bits(g.want)
		gotBits := math.Float64bits(got)
		if wantBits != gotBits {
			t.Errorf("0x%04X: number mismatch; expected 0x%04X (%v), got 0x%04X (%v)", g.bits, wantBits, g.want, gotBits, got)
		}
	}
}

func TestNewFromFloat64(t *testing.T) {
	golden := []struct {
		in    float64
		exact bool
		want  uint16
	}{
		// Special numbers.

		// +NaN
		/*
			{in: math.NaN(), exact: true, want: 0x7E00},
			// -NaN
			{in: -math.NaN(), exact: true, want: 0xFE00},
			// +Inf
			{in: math.Inf(1), exact: true, want: 0x7C00},
			// -Inf
			{in: math.Inf(-1), exact: true, want: 0xFC00},
			// +0
			{in: 0.0, exact: true, want: 0x0000},
			// -0
			{in: math.Copysign(0.0, -1), exact: true, want: 0x8000},
		*/

		// From https://reviews.llvm.org/rL237161

		// Normalized numbers.
		{in: 0.5, exact: true, want: 0x3800},
		{in: -0.5, exact: true, want: 0xB800},
		{in: 1.5, exact: true, want: 0x3E00},
		{in: -1.5, exact: true, want: 0xBE00},
		{in: 2.5, exact: true, want: 0x4100},
		{in: -2.5, exact: true, want: 0xC100},

		// Denormalized numbers.
		{in: float64FromString("0x1.0p-20"), exact: true, want: 0x0010},
		{in: float64FromString("0x1.0p-24"), exact: true, want: 0x0001},
		{in: float64FromString("-0x1.0p-24"), exact: true, want: 0x8001},
		//{in: float64FromString("0x1.5p-25"), exact: true, want: 0x0001},

		// Rounding.
		// TODO: Handle rounding.
		//{in: 3.14, exact: true, want: 0x4248},
		//{in: -3.14, exact: true, want: 0xC248},
		//{in: 3.1415926535, exact: true, want: 0x4248},
		//{in: -3.1415926535, exact: true, want: 0xC248},
		//{in: float64FromString("0x1.987124876876324p+100"), exact: true, want: 0x7C00},
		{in: float64FromString("0x1.988p+12"), exact: true, want: 0x6E62},
		{in: float64FromString("0x1.0p+0"), exact: true, want: 0x3C00},
		{in: float64FromString("0x1.0p-14"), exact: true, want: 0x0400},
		// rounded to zero
		//{in: float64FromString("0x1.0p-25"), exact: true, want: 0x0000},
		//{in: float64FromString("-0x1.0p-25"), exact: true, want: 0x8000},
		// max (precise)
		{in: 65504.0, exact: true, want: 0x7BFF},
	}

	for _, g := range golden {
		f, exact := NewFromFloat64(g.in)
		if g.exact != exact {
			t.Errorf("%v: exact mismatch; expected %v, got %v", g.in, g.exact, exact)
		}
		got := f.Bits()
		if g.want != got {
			t.Errorf("%v: number mismatch; expected 0x%04X, got 0x%04X", g.in, g.want, got)
		}
	}
}

func TestNewFromFloat32(t *testing.T) {
	golden := []struct {
		uint32Float uint32
		a           uint16
		str         string
	}{
		// Special numbers.

		// +NaN
		{uint32Float: 0x7F800000, a: 0x7C00, str: "+Nan not equal"},
		// -NaN
		{uint32Float: 0xFF800000, a: 0xFC00, str: "-Nan not equal"},
		// +Inf
		{uint32Float: 0x7FC00000, a: 0x7E00, str: "+Inf not equal"},
		// -Inf
		{uint32Float: 0xFFC00000, a: 0xFE00, str: "-Inf not equal"},
		// +0
		{uint32Float: 0x00000000, a: 0, str: "+0 not equal"},
		// -0
		{uint32Float: 0x80000000, a: 0x8000, str: "-0 not equal"},

		// Normalized numbers.
		{uint32Float: math.Float32bits(0.5), a: 0x3800, str: "+0.5 not equal"},
		{uint32Float: math.Float32bits(-0.5), a: 0xB800, str: "-0.5 not equal"},
		{uint32Float: math.Float32bits(1.5), a: 0x3E00, str: "+1.5 not equal"},
		{uint32Float: math.Float32bits(-1.5), a: 0xBE00, str: "-1.5 not equal"},
		{uint32Float: math.Float32bits(2.5), a: 0x4100, str: "+2.5 not equal"},
		{uint32Float: math.Float32bits(-2.5), a: 0xC100, str: "-2.5 not equal"},
	}

	for _, g := range golden {
		f, _ := NewFromFloat32(math.Float32frombits(g.uint32Float))
		a := f.Bits()
		assert.Equal(t, g.a, a, g.str)
	}
}

func TestNewFromFloat64(t *testing.T) {
	golden := []struct {
		uint64Float uint64
		a           uint16
		str         string
	}{
		// Special numbers.

		// +NaN
		{uint64Float: 0x7FF0000000000000, a: 0x7C00, str: "+Nan not equal"},
		// -NaN
		{uint64Float: 0xFFF0000000000000, a: 0xFC00, str: "-Nan not equal"},
		// +Inf
		{uint64Float: 0x7FF8000000000000, a: 0x7E00, str: "+Inf not equal"},
		// -Inf
		{uint64Float: 0xFFF8000000000000, a: 0xFE00, str: "-Inf not equal"},
		// +0
		{uint64Float: 0x00000000, a: 0, str: "+0 not equal"},
		// -0
		{uint64Float: 0x8000000000000000, a: 0x8000, str: "-0 not equal"},

		// Normalized numbers.
		{uint64Float: math.Float64bits(0.5), a: 0x3800, str: "+0.5 not equal"},
		{uint64Float: math.Float64bits(-0.5), a: 0xB800, str: "-0.5 not equal"},
		{uint64Float: math.Float64bits(1.5), a: 0x3E00, str: "+1.5 not equal"},
		{uint64Float: math.Float64bits(-1.5), a: 0xBE00, str: "-1.5 not equal"},
		{uint64Float: math.Float64bits(2.5), a: 0x4100, str: "+2.5 not equal"},
		{uint64Float: math.Float64bits(-2.5), a: 0xC100, str: "-2.5 not equal"},
	}

	for _, g := range golden {
		f, _ := NewFromFloat64(math.Float64frombits(g.uint64Float))
		a := f.Bits()
		assert.Equal(t, g.a, a, g.str)
	}
}

func float64FromString(s string) float64 {
	x, _, err := big.ParseFloat(s, 0, 53, big.ToNearestEven)
	if err != nil {
		panic(err)
	}
	// TODO: Check accuracy?
	y, _ := x.Float64()
	return y
}
