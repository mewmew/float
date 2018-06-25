// Package binary16 implements encoding and decoding of IEEE 754 half precision
// floating-point numbers.
//
// https://en.wikipedia.org/wiki/Half-precision_floating-point_format
package binary16

// Float is a floating-point number in IEEE 754 half precision format.
type Float struct {
	// Sign, exponent and fraction.
	//
	//    1 bit:   sign
	//    5 bits:  exponent
	//    10 bits: fraction
	a uint16
}

// NewFromBits returns the floating-point number corresponding to the IEEE 754
// half precision binary representation.
func NewFromBits(a uint16) Float {
	return Float{a: a}
}

// NewFromFloat32 returns the nearest half precision floating-point number for x
// and a bool indicating whether f represents x exactly.
func NewFromFloat32(x float32) (f Float, exact bool) {
	panic("not yet implemented")
}

// NewFromFloat64 returns the nearest half precision floating-point number for x
// and a bool indicating whether f represents x exactly.
func NewFromFloat64(x float64) (f Float, exact bool) {
	panic("not yet implemented")
}

// Bits returns the IEEE 754 half precision binary representation of f.
func (f Float) Bits() uint16 {
	return f.a
}

// Float32 returns the float32 representation of f.
func (f Float) Float32() float32 {
	panic("not yet implemented")
}

// Float64 returns the float64 representation of f.
func (f Float) Float64() float64 {
	panic("not yet implemented")
}
