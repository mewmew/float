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

// Positive and negative Not-a-Number, infinity and zero.
var (
	// +NaN
	NaN = Float{high: math.NaN(), low: 0}
	// -NaN
	NegNaN = Float{high: -math.NaN(), low: 0}
	// +Inf
	Inf = Float{high: math.Inf(1), low: 0}
	// -Inf
	NegInf = Float{high: -math.Inf(-1), low: 0}
	// +zero
	Zero = Float{high: 0, low: 0}
	// -zero
	NegZero = Float{high: math.Copysign(0, -1), low: 0}
)

// Float is a floating-point number in double-double format.
type Float struct {
	// where a long double value is regarded as the exact sum of two double-precision values, giving at least a 106-bit precision
	high float64
	low  float64
}

// NewFromBits returns the floating-point number corresponding to the
// double-double representation.
func NewFromBits(a, b uint64) Float {
	return Float{high: math.Float64frombits(a), low: math.Float64frombits(b)}
}

// NewFromFloat32 returns the nearest double-double precision floating-point
// number for x and the accuracy of the conversion.
func NewFromFloat32(x float32) (f Float, exact big.Accuracy) {
	f, acc := NewFromFloat64(float64(x))
	if acc == big.Exact {
		_, acc = f.Float32()
	}
	return f, acc
}

// NewFromFloat64 returns the nearest double-double precision floating-point
// number for x and the accuracy of the conversion.
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
	r := Float{high: x, low: 0}
	br, _ := r.Big()
	return r, br.Acc()
}

// NewFromBig returns the nearest double-double floating-point number for x and
// the accuracy of the conversion.
func NewFromBig(x *big.Float) (Float, big.Accuracy) {
	// +-Inf
	zero := big.NewFloat(0).SetPrec(precision)
	switch {
	case x.IsInf():
		if x.Signbit() {
			// -Inf
			return NegInf, big.Exact
		}
		// +Inf
		return Inf, big.Exact
	// +-zero
	case x.Cmp(zero) == 0:
		if x.Signbit() {
			// -zero
			return NegZero, big.Exact
		}
		// +zero
		return Zero, big.Exact
	}

	high, acc := x.Float64()
	// TODO: figure out how to implement NewFromBig in a good way. Currently all
	// added precision of low is lost. A proper implementation would store half of
	// the mantissa in high and half in low, adjusting the exponent such that
	//
	//    high + low = f
	//
	// and
	//
	//    |high| >= |low|
	return Float{high: high, low: 0}, acc
}

// Bits returns the double-double binary representation of f.
func (f Float) Bits() (a, b uint64) {
	return math.Float64bits(f.high), math.Float64bits(f.low)
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
	if f.IsNaN() {
		return x, true
	}
	x.Add(big.NewFloat(f.high), big.NewFloat(f.low))
	return x, false
}

// IsNaN returns true if the Float is NaN
func (f Float) IsNaN() bool {
	// NaN + NaN should be NaN in consideration
	return math.IsNaN(f.high) || math.IsNaN(f.low)
}
