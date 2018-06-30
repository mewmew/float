// Package binary128 implements encoding and decoding of IEEE 754 quadruple
// precision floating-point numbers.
//
// https://en.wikipedia.org/wiki/Quadruple-precision_floating-point_format
package binary128

import "math"

// Float is a floating-point number in IEEE 754 quadruple precision format.
type Float struct {
	// Sign, exponent and fraction.
	//
	//    1 bit:    sign
	//    15 bits:  exponent
	//    112 bits: fraction
	a uint64
	b uint64
}

// NewFromBits returns the floating-point number corresponding to the IEEE 754
// quadruple precision binary representation.
func NewFromBits(a, b uint64) Float {
	return Float{a: a, b: b}
}

// NewFromFloat32 returns the nearest quadruple precision floating-point number
// for x and a bool indicating whether f represents x exactly.
func NewFromFloat32(x float32) (f Float, exact bool) {
	intRep := math.Float32bits(x)
	sign := intRep & 0x80000000 != 0
	mant := intRep & 0x7fffff
	exp := intRep & 0x7f800000 >> 23

	switch exp {
	// 0b11111111
	case 0xFF:
		// NaN or Inf
		var a uint64
		if mant == 0 {
			// +-Inf
			a = 0x7FFF000000000000
			if sign {
				a = 0xFFFF000000000000
			}
			return Float{a:a, b:0}, true
		}
		// +-NaN

		a = 0
		if sign {
			a = 0x8000000000000000
		}
		a = a | 0x7FFF000000000000

		newMant := uint64(mant) << uint64(25)
		a = a | newMant

		return Float{a:a, b:0}, true
		// 0b00000000
	case 0x00:
		if mant == 0 {
			// +-Zero
			var a uint64
			a = 0x0000000000000000
			if sign {
				a = 0x8000000000000000
			}
			return Float{a:a, b:0}, true
		}
	}

	var a uint64 = 0
	if sign {
		a = 0x8000000000000000
	}

	var newExp uint64
	newExp = uint64(exp - 127 + 1023) << uint64(52)
	a = a | newExp

	var newMant uint64
	newMant = uint64(mant) << uint64(41)

	a = a | newMant

	return Float{a: a, b: 0}, true
}

// NewFromFloat64 returns the nearest quadruple precision floating-point number
// for x and a bool indicating whether f represents x exactly.
func NewFromFloat64(x float64) (f Float, exact bool) {
	panic("not yet implemented")
}

// Bits returns the IEEE 754 quadruple precision binary representation of f.
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
