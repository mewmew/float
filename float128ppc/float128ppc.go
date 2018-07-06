// Package float128ppc implements encoding and decoding of IBM extended double
// floating-point numbers.
//
// https://en.wikipedia.org/wiki/IBM_Floating_Point_Architecture#Extended-precision_128-bit
package float128ppc

// Float is a floating-point number in IBM extended double format.
type Float struct {
	// Sign, exponent and fraction.
	//
	//    1 bit:    sign
	//    7 bits:   exponent
	//    112 bits: fraction
	//    8 bits:   unused
	a uint64
	b uint64
}

// NewFromBits returns the floating-point number corresponding to the IBM
// extended double representation.
func NewFromBits(a uint16, b uint64) Float {
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
