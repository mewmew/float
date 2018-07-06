// Package float80x86 implements encoding and decoding of x86 extended precision
// floating-point numbers.
//
// https://en.wikipedia.org/wiki/Extended_precision#x86_extended_precision_format
package float80x86

// Float is a floating-point number in x86 extended precision format.
type Float struct {
	// Sign and exponent.
	//
	//    1 bit:   sign
	//    15 bits: exponent
	se uint16
	// Integer part and fraction.
	//
	//    1 bit:   integer part
	//    63 bits: fraction
	m uint64
}

// NewFromBits returns the floating-point number corresponding to the x86
// extended precision representation.
func NewFromBits(se uint16, m uint64) Float {
	return Float{se: se, m: m}
}

// NewFromFloat32 returns the nearest x86 extended precision floating-point
// number for x and a bool indicating whether f represents x exactly.
func NewFromFloat32(x float32) (f Float, exact bool) {
	panic("not yet implemented")
}

// NewFromFloat64 returns the nearest x86 extended precision floating-point
// number for x and a bool indicating whether f represents x exactly.
func NewFromFloat64(x float64) (f Float, exact bool) {
	panic("not yet implemented")
}

// Bits returns the x86 extended precision binary representation of f.
func (f Float) Bits() (se uint16, m uint64) {
	return f.se, f.m
}

// Float32 returns the float32 representation of f.
func (f Float) Float32() float32 {
	panic("not yet implemented")
}

// Float64 returns the float64 representation of f.
func (f Float) Float64() float64 {
	panic("not yet implemented")
}
