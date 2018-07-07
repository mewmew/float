// Package binary16 implements encoding and decoding of IEEE 754 half precision
// floating-point numbers.
//
// https://en.wikipedia.org/wiki/Half-precision_floating-point_format
package binary16

import (
	"fmt"
	"math"
	"math/big"
)

const (
	// precision specifies the number of bits in the mantissa (including the
	// implicit lead bit).
	precision = 11
	// exponent bias.
	bias = 15
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
	y := float64(x)
	f, exact = NewFromFloat64(y)
	return f, exact && x == float32(y)
}

// NewFromFloat64 returns the nearest half precision floating-point number for x
// and a bool indicating whether f represents x exactly.
func NewFromFloat64(x float64) (f Float, exact bool) {
	// +-NaN
	switch {
	case math.IsNaN(x):
		if math.Signbit(x) {
			return Float{bits: 0xFE00}, true
		}
		return Float{bits: 0x7E00}, true
	}
	y := big.NewFloat(x)
	y.SetPrec(precision)
	y.SetMode(big.ToNearestEven)
	// TODO: check accuracy after setting precision?
	return NewFromBig(y)
}

// NewFromBig returns the nearest half precision floating-point number for x and
// a bool indicating whether f represents x exactly.
func NewFromBig(x *big.Float) (f Float, exact bool) {
	fmt.Println("x:", x)
	zero := &big.Float{}
	switch {
	// +-Inf
	case x.IsInf():
		if x.Signbit() {
			// -Inf
			return Float{bits: 0xFC00}, true
		}
		// +Inf
		return Float{bits: 0x7C00}, true
	// +-zero
	case x.Cmp(zero) == 0:
		// -zero
		if x.Signbit() {
			return Float{bits: 0x8000}, true
		}
		// +zero
		return Float{bits: 0x0000}, true
	}

	// Sign
	var bits uint16
	if x.Signbit() {
		bits |= 0x8000
	}
	fmt.Println("sign:", x.Signbit())
	fmt.Printf("bits: %04X\n", bits)

	// Exponent and mantissa.
	mant := &big.Float{}
	exponent := x.MantExp(mant)
	// TODO: Check overflow and set exact.
	fmt.Println("exponent:", exponent)
	exp := exponent - 1 + bias
	// 0b11111
	exact = (exp &^ 0x1F) == 0
	bits |= uint16(exp&0x1F) << 10
	fmt.Println("exp:", exp)
	fmt.Printf("bits: %04X\n", bits)

	fmt.Println("mant:", mant)
	if mant.Signbit() {
		mant.Neg(mant)
	}
	mant.SetMantExp(mant, precision)
	fmt.Println("mant:", mant)
	if !mant.IsInt() {
		// TODO: handle truncation.
		panic("truncation")
	}
	mantissa, _ := mant.Uint64()
	fmt.Println("mantissa:", mantissa)
	mantissa &^= 1024 // clear implicit lead bit.
	fmt.Println("mantissa:", mantissa)

	//m, _ := mant.Rat(nil)
	//fmt.Println("m:", m)

	// 0b11111111111 (including implicit lead bit)
	exact = exact && (mantissa&^0x7FF) == 0
	mantissa &= 0x7FF
	bits |= uint16(mantissa)
	fmt.Println("mantissa:", mantissa)
	fmt.Printf("bits: %04X\n", bits)
	fmt.Println()
	return Float{bits: bits}, exact
}

// Bits returns the IEEE 754 half precision binary representation of f.
func (f Float) Bits() uint16 {
	return f.bits
}

// Float32 returns the float32 representation of f.
func (f Float) Float32() float32 {
	x, nan := f.Big()
	if nan {
		if x.Signbit() {
			return float32(-math.NaN())
		}
		return float32(math.NaN())
	}
	y, _ := x.Float32()
	return y
}

// Float64 returns the float64 representation of f.
func (f Float) Float64() float64 {
	x, nan := f.Big()
	if nan {
		if x.Signbit() {
			return -math.NaN()
		}
		return math.NaN()
	}
	y, _ := x.Float64()
	return y
}

// Big returns the multi-precision floating-point number representation of f and
// a boolean indicating whether f is Not-a-Number.
func (f Float) Big() (x *big.Float, nan bool) {
	signbit := f.signbit()
	exp := f.exp()
	mant := f.mant()
	x = big.NewFloat(0)
	x.SetPrec(precision)
	x.SetMode(big.ToNearestEven)

	// ref: https://en.wikipedia.org/wiki/Half-precision_floating-point_format#Exponent_encoding
	//
	// 0b00001 - 0b11110
	// Normalized number.
	//
	//    (-1)^signbit * 2^(exp-15) * 1.mant_2
	lead := 1
	exponent := int(exp) - bias

	switch exp {
	// 0b11111
	case 0x1F:
		// Inf or NaN
		if mant == 0 {
			// +-Inf
			x.SetInf(signbit)
			return x, false
		}
		// +-NaN
		if signbit {
			x.Neg(x)
		}
		return x, true
	// 0b00000
	case 0x00:
		if mant == 0 {
			// +-Zero
			if signbit {
				x.Neg(x)
			}
			return x, false
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
		panic(err)
	}
	return x, false
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
