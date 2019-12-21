// Package float128ppc implements encoding and decoding of double-double
// floating-point numbers.
//
// https://en.wikipedia.org/wiki/Quadruple-precision_floating-point_format#Double-double_arithmetic
package float128ppc

import (
	"math"
	"math/big"
)

const (
	// precision specifies the number of bits in the mantissa (including the
	// implicit lead bit).
	precision = 106
)

var (
	NegNaN = Float{a: -math.NaN(), b: 0}
	NaN    = Float{a: math.NaN(), b: 0}
)

// Float is a floating-point number in double-double format.
type Float struct {
	// where a long double value is regarded as the exact sum of two double-precision values, giving at least a 106-bit precision
	a float64
	b float64
}

// NewFromBits returns the floating-point number corresponding to the IBM
// extended double representation.
func NewFromBits(a, b uint64) Float {
	return Float{a: math.Float64frombits(a), b: math.Float64frombits(b)}
}

// NewFromFloat32 returns the nearest IBM extended double floating-point number
// for x and a bool indicating whether f represents x exactly.
func NewFromFloat32(x float32) (f Float, exact big.Accuracy) {
	f, acc := NewFromFloat64(float64(x))
	if acc == big.Exact {
		_, acc = f.Float32()
	}
	return f, acc
}

// NewFromFloat64 returns the nearest IBM extended double floating-point number
// for x and a bool indicating whether f represents x exactly.
func NewFromFloat64(x float64) (f Float, exact big.Accuracy) {
	// +-NaN
	switch {
	case math.IsNaN(x):
		if math.Signbit(x) {
			// -NaN
			return NegNaN, big.Exact
		}
		// +NaN
		return NaN, big.Exact
	}
	r := Float{a: x, b: 0}
	br, _ := r.Big()
	return r, br.Acc()
}

// Bits returns the IBM extended double binary representation of f.
func (f Float) Bits() (a, b uint64) {
	return math.Float64bits(f.a), math.Float64bits(f.b)
}

// Float32 returns the float32 representation of f.
func (f Float) Float32() (float32, big.Accuracy) {
	x, nan := f.Big()
	if nan {
		if x.Signbit() {
			return float32(-math.NaN()), big.Exact
		}
		return float32(math.NaN()), big.Exact
	}
	return x.Float32()
}

// Float64 returns the float64 representation of f.
func (f Float) Float64() (float64, big.Accuracy) {
	x, nan := f.Big()
	if nan {
		if x.Signbit() {
			return -math.NaN(), big.Exact
		}
		return math.NaN(), big.Exact
	}
	return x.Float64()
}

// Big returns the multi-precision floating-point number representation of f and
// a boolean indicating whether f is Not-a-Number.
func (f Float) Big() (x *big.Float, nan bool) {
	x = big.NewFloat(0)
	x.SetPrec(precision)
	x.SetMode(big.ToNearestEven)
	if math.IsNaN(f.a) || math.IsNaN(f.b) {
		return x, true
	}
	h := big.NewFloat(f.a)
	l := big.NewFloat(f.b)
	x.Add(h, l)
	return x, false
}
