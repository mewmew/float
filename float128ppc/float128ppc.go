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

// Float is a floating-point number in double-double format.
type Float struct {
	// where a long double value is regarded as the exact sum of two double-precision values, giving at least a 106-bit precision
	a uint64
	b uint64
}

// NewFromBits returns the floating-point number corresponding to the IBM
// extended double representation.
func NewFromBits(a, b uint64) Float {
	return Float{a: a, b: b}
}

// NewFromFloat32 returns the nearest IBM extended double floating-point number
// for x and a bool indicating whether f represents x exactly.
func NewFromFloat32(x float32) (f Float, exact bool) {
	panic("not yet implemented")
}

// NewFromFloat64 returns the nearest IBM extended double floating-point number
// for x and a bool indicating whether f represents x exactly.
func NewFromFloat64(x float64) (f Float, exact bool) {
	panic("not yet implemented")
}

// Bits returns the IBM extended double binary representation of f.
func (f Float) Bits() (a, b uint64) {
	return f.a, f.b
}

// Float32 returns the float32 representation of f.
func (f Float) Float32() float32 {
	panic("not yet implemented")
}

// Float64 returns the float64 representation of f.
func (f Float) Float64() float64 {
	panic("not yet implemented")
}

// Big returns the multi-precision floating-point number representation of f and
// a boolean indicating whether f is Not-a-Number.
func (f Float) Big() (x *big.Float, nan bool) {
	x = big.NewFloat(0)
	x.SetPrec(precision)
	x.SetMode(big.ToNearestEven)
	a := math.Float64frombits(f.a)
	if math.IsNaN(a) {
		return x, true
	}
	b := math.Float64frombits(f.b)
	if math.IsNaN(b) {
		return x, true
	}
	h := big.NewFloat(a)
	l := big.NewFloat(b)
	x.Add(h, l)
	return x, false
}
