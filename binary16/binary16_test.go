package binary16_test

import (
	"fmt"
	"math"
	"math/big"
	"testing"

	"github.com/mewmew/floats/binary16"
)

func TestNewFromBits(t *testing.T) {
	golden := []struct {
		bits uint16
		want float64
	}{
		//{bits: 0x7E00, want: +math.NaN()},
		//{bits: 0xFE00, want: -math.NaN()},
		//{bits: 0x7C00, want: math.Inf(1)},
		//{bits: 0xFC00, want: math.Inf(-1)},
		//{bits: 0x0000, want: +0.0},
		//{bits: 0x8000, want: math.Copysign(0.0, -1)},
		//{bits: 0x3800, want: +0.5},
		//{bits: 0xB800, want: -0.5},
		{bits: 0x3E00, want: +1.5},
		{bits: 0xBE00, want: -1.5},
		{bits: 0x4100, want: +2.5},
		{bits: 0xC100, want: -2.5},
		{bits: 0x4248, want: +3.14},
		{bits: 0xC248, want: -3.14},

		// From https://reviews.llvm.org/rL237161

		// NaN
		{bits: 0x7e00, want: math.NaN()},
		// inf
		{bits: 0x7c00, want: math.Inf(1)},
		{bits: 0xfc00, want: math.Inf(-1)},
		// zero
		{bits: 0x0000, want: +0.0},
		{bits: 0x8000, want: math.Copysign(0.0, -1)},

		{bits: 0x4248, want: 3.1415926535},
		{bits: 0xc248, want: -3.1415926535},
		{bits: 0x7c00, want: float64FromString("0x1.987124876876324p+100")},
		{bits: 0x6e62, want: float64FromString("0x1.988p+12")},
		{bits: 0x3c00, want: float64FromString("0x1.0p+0")},
		{bits: 0x0400, want: float64FromString("0x1.0p-14")},
		// denormal
		{bits: 0x0010, want: float64FromString("0x1.0p-20")},
		{bits: 0x0001, want: float64FromString("0x1.0p-24")},
		{bits: 0x8001, want: float64FromString("-0x1.0p-24")},
		{bits: 0x0001, want: float64FromString("0x1.5p-25")},
		// and back to zero
		{bits: 0x0000, want: float64FromString("0x1.0p-25")},
		{bits: 0x8000, want: float64FromString("-0x1.0p-25")},
		// max (precise)
		{bits: 0x7bff, want: 65504.0},
		// max (rounded)
		{bits: 0x7bff, want: 65504.0},
		// max (to +inf)
		{bits: 0x7c00, want: math.Inf(1)},
		{bits: 0xfc00, want: math.Inf(-1)},
	}

	for _, g := range golden {
		f := binary16.NewFromBits(g.bits)
		got := f.Float64()
		fmt.Printf("0x%04X\n", g.bits)
		fmt.Println("   got:  ", got)
		fmt.Println("   want: ", g.want)
		fmt.Println()
	}
}

func float64FromString(s string) float64 {
	x, _, err := big.ParseFloat(s, 0, 53, big.ToNearestEven)
	if err != nil {
		panic(err)
	}
	//x.ParseFloat()
	//x, err := strconv.ParseFloat(s, 64)
	//if err != nil {
	//	panic(err)
	//}
	y, acc := x.Float64()
	fmt.Println("acc:", acc)
	return y
}
