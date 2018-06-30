package binary128

import (
	"testing"
	"math"
	"math/big"
	"github.com/stretchr/testify/assert"
)

func TestNewFromFloat32(t *testing.T) {
	golden := []struct {
		uint32Float uint32
		a uint64
		b uint64
		str string
	}{
		// Special numbers.

		// +NaN
		{uint32Float: 0x7F800000, a: 0x7FFF000000000000, b: 0, str: "+Nan not equal"},
		// -NaN
		{uint32Float: 0xFF800000, a: 0xFFFF000000000000, b: 0, str: "-Nan not equal"},
		// +Inf
		{uint32Float: 0x7FC00000, a: 0x7FFF800000000000, b: 0, str: "+Inf not equal"},
		// -Inf
		{uint32Float: 0xFFC00000, a: 0xFFFF800000000000, b: 0, str: "-Inf not equal"},
		// +0
		{uint32Float: 0x00000000, a: 0, b: 0, str: "+0 not equal"},
		// -0
		{uint32Float: 0x80000000, a: 0x8000000000000000, b :0, str: "-0 not equal"},


	}

	for _, g := range golden {
		f, _ := NewFromFloat32(math.Float32frombits(g.uint32Float))
		a, b := f.Bits()
		assert.Equal(t, g.a, a, g.str)
		assert.Equal(t, g.b, b, g.str)
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