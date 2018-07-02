// Package floats implements multi-precision floating-point numbers.
package floats

import (
	"math"
	"math/big"
)

// A Float is a multi-precision floating-point number.
type Float struct {
	*big.Float
	// NaN denotes not-a-number.
	NaN bool
}

// New returns a zero value multi-precision floating-point number.
func New() *Float {
	return &Float{
		Float: &big.Float{},
	}
}

// Float64 returns the float64 value nearest to f and the associated accuracy.
func (f *Float) Float64() (float64, big.Accuracy) {
	if f.NaN {
		if f.Signbit() {
			return -math.NaN(), big.Exact
		}
		return math.NaN(), big.Exact
	}
	return f.Float.Float64()
}
