// Package binary16 implements encoding and decoding of IEEE 754 half precision
// floating-point numbers.
//
// https://en.wikipedia.org/wiki/Half-precision_floating-point_format
package binary16

import (
	"fmt"
	"log"
	"math/big"

	"github.com/mewmew/floats"
)

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
	x := f.big()
	// TODO: Check accuracy?
	y, _ := x.Float64()
	return y
}

// big returns the multi-precision floating-point number representation of f.
func (f Float) big() *floats.Float {
	signbit := f.signbit()
	exp := f.exp()
	mant := f.mant()
	x := floats.New()
	x.SetPrec(11)
	x.SetMode(big.ToNearestEven)
	// ref: https://en.wikipedia.org/wiki/Half-precision_floating-point_format#Exponent_encoding

	// 0b00001 - 0b11110
	// Normalized number.
	//
	//    (-1)^signbit * 2^(exp-15) * 1.mant_2
	lead := 1
	const bias = 15
	exponent := int(exp) - bias

	switch exp {
	// 0b11111
	case 0x1F:
		// NaN or Inf
		if mant == 0 {
			// +-Inf
			x.SetInf(signbit)
			return x
		}
		// +-NaN
		x.NaN = true
		if signbit {
			x.Neg(x.Float)
		}
		return x
	// 0b00000
	case 0x00:
		if mant == 0 {
			// +-Zero
			if signbit {
				x.Neg(x.Float)
			}
			return x
		}
		// Denormalized number.
		//
		//    (-1)^signbit * 2^(-14) * 0.mant_2
		lead = 0
		exponent = -14
	}

	// number = [ sign ] [ prefix ] mantissa [ exponent ] | infinity .
	sign := "+"
	if signbit {
		sign = "-"
	}
	s := fmt.Sprintf("%s0b%d.%010bp%d", sign, lead, mant, exponent)
	_, _, err := x.Parse(s, 0)
	if err != nil {
		log.Fatalf("%+v", err)
	}
	return x
}

// signbit reports whether f is negative or negative 0.
func (f Float) signbit() bool {
	// 0b1000000000000000
	return f.a&0x8000 != 0
}

// exp returns the exponent of f.
func (f Float) exp() uint16 {
	// 0b0111110000000000
	return f.a & 0x7C00 >> 10
}

// mant returns the mantissa of f.
func (f Float) mant() uint16 {
	// 0b0000001111111111
	return f.a & 0x03FF
}
