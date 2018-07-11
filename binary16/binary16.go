// Package binary16 implements encoding and decoding of IEEE 754 half precision
// floating-point numbers.
//
// https://en.wikipedia.org/wiki/Half-precision_floating-point_format
package binary16

import (
	"fmt"
	"log"
	"math"
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
	bits uint16
}

// NewFromBits returns the floating-point number corresponding to the IEEE 754
// half precision binary representation.
func NewFromBits(bits uint16) Float {
	return Float{bits: bits}
}

// NewFromFloat32 returns the nearest half precision floating-point number for x
// and a bool indicating whether f represents x exactly.
func NewFromFloat32(x float32) (f Float, exact bool) {
	intRep := math.Float32bits(x)
	sign := intRep&0x80000000 != 0
	mant := intRep & 0x7fffff
	exp := int32(intRep & 0x7f800000 >> 23)

	var mantTruncMask uint32 = 0x7FE000
	var mantShift uint32 = 13

	switch exp {
	// 0b11111
	case 0xFF:
		// NaN or Inf
		var a uint16
		if mant == 0 {
			// +-Inf
			a = 0x7C00
			if sign {
				a = 0xFC00
			}
			return Float{bits: a}, true
		}
		// +-NaN

		a = 0
		if sign {
			a = 0x8000
		}
		a = a | 0x7C00

		truncMant := mant & mantTruncMask
		exact := true
		if mant-truncMant > 0 {
			exact = false
		}

		newMant := uint16(truncMant >> mantShift)
		a = a | newMant

		return Float{bits: a}, exact
		// 0b00000000
	case 0x00:
		if mant == 0 {
			// +-Zero
			var a uint16
			if sign {
				a = 0x8000
			}
			return Float{bits: a}, true
		}
	}

	var a uint16
	if sign {
		a = 0x8000
	}

	expVal := exp - 127
	if expVal > 16 { // Inf with not exact
		a |= 0x7C00
		return Float{bits: a}, false
	}

	if expVal < -16 { // Zero with not exact
		return Float{bits: a}, false
	}

	truncMant := mant & mantTruncMask
	exact = true
	if mant-truncMant > 0 {
		exact = false
	}

	newMant := uint16(truncMant >> mantShift)
	a = a | newMant

	newExp := uint16(expVal+15) << 10
	a |= newExp

	return Float{bits: a}, exact
}

// NewFromFloat64 returns the nearest half precision floating-point number for x
// and a bool indicating whether f represents x exactly.
func NewFromFloat64(x float64) (f Float, exact bool) {
	intRep := math.Float64bits(x)
	sign := intRep&0x8000000000000000 != 0
	exp := int64(intRep & 0x7FF0000000000000 >> 52)
	mant := intRep & 0xFFFFFFFFFFFFF

	var mantTruncMask uint64 = 0xFFC0000000000
	var mantShift uint64 = 42

	switch exp {
	// 0b11111
	case 0x7FF:
		// NaN or Inf
		var a uint16
		if mant == 0 {
			// +-Inf
			a = 0x7C00
			if sign {
				a = 0xFC00
			}
			return Float{bits: a}, true
		}
		// +-NaN

		a = 0
		if sign {
			a = 0x8000
		}
		a = a | 0x7C00

		truncMant := mant & mantTruncMask
		exact := true
		if mant-truncMant > 0 {
			exact = false
		}

		newMant := uint16(truncMant >> mantShift)
		a = a | newMant

		return Float{bits: a}, exact
		// 0b00000000
	case 0x00:
		if mant == 0 {
			// +-Zero
			var a uint16
			if sign {
				a = 0x8000
			}
			return Float{bits: a}, true
		}
	}

	var a uint16
	if sign {
		a = 0x8000
	}

	exp = exp - 1023
	if exp > 16 { // Inf with not exact
		a |= 0x7C00
		return Float{bits: a}, false
	}

	if exp < -16 { // Zero with not exact
		return Float{bits: a}, false
	}

	truncMant := mant & mantTruncMask
	exact = true
	if mant-truncMant > 0 {
		exact = false
	}

	newMant := uint16(truncMant >> mantShift)
	a = a | newMant

	newExp := uint16(exp+15) << 10
	a |= newExp

	return Float{bits: a}, exact
}

// Bits returns the IEEE 754 half precision binary representation of f.
func (f Float) Bits() uint16 {
	return f.bits
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
	return f.bits&0x8000 != 0
}

// exp returns the exponent of f.
func (f Float) exp() uint16 {
	// 0b0111110000000000
	return f.bits & 0x7C00 >> 10
}

// mant returns the mantissa of f.
func (f Float) mant() uint16 {
	// 0b0000001111111111
	return f.bits & 0x03FF
}
