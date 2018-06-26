package floats

import (
	"math"
	"math/big"
)

type Float struct {
	*big.Float
	NaN bool
}

func New() *Float {
	return &Float{
		Float: &big.Float{},
	}
}

func (f *Float) Float64() (float64, big.Accuracy) {
	if f.NaN {
		return math.NaN(), big.Exact
	}
	return f.Float.Float64()
}
