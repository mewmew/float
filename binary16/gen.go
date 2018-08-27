//+build ignore

package main

import (
	"flag"
	"fmt"
	"html/template"
	"log"
	"math"
	"math/big"
	"os"

	"github.com/pkg/errors"
)

func main() {
	var out string
	flag.StringVar(&out, "o", "extra_test.go", "test cases output path")
	flag.Parse()
	if err := dumpTest(out); err != nil {
		log.Fatalf("%+v", err)
	}
}

func dumpTest(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return errors.WithStack(err)
	}
	defer f.Close()
	t, err := template.ParseFiles("extra_test.tmpl")
	if err != nil {
		return errors.WithStack(err)
	}
	data := map[string][]string{
		"normalized":   getNormalized(),
		"denormalized": getDenormalized(),
	}
	if err := t.Execute(f, data); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

const bias = 15

func getNormalized() []string {
	var ns []string
	// normalized
	// 0b00001 - 0b11110
	//
	//    (-1)^signbit * 2^(exp-15) * 1.mant_2
	const lead = 1
	for signbit := 0; signbit <= 1; signbit++ {
		for exp := 1; exp <= 0x1E; exp++ {
			// 0b1111111111
			for mant := 0; mant <= 0x3FF; mant++ {
				s := fmt.Sprintf("%s0b%d.%010bp0", "+", lead, mant)
				m, _, err := big.ParseFloat(s, 0, 53, big.ToNearestEven)
				if err != nil {
					panic(err)
				}
				mantissa, acc := m.Float64()
				if acc != big.Exact {
					panic("not exact")
				}
				x := math.Pow(-1, float64(signbit)) * math.Pow(2, float64(exp)-bias) * mantissa
				bits := uint16(signbit) << 15
				bits |= uint16(exp) << 10
				bits |= uint16(mant)
				n := fmt.Sprintf("{bits: 0x%04X, want: %v}", bits, x)
				ns = append(ns, n)
			}
		}
	}
	return ns
}

func getDenormalized() []string {
	var ds []string
	// denormalized
	// 0b00000
	//
	//    (-1)^signbit * 2^(-14) * 0.mant_2
	const lead = 0
	for signbit := 0; signbit <= 1; signbit++ {
		// 0b1111111111
		const exp = 0
		for mant := 0; mant <= 0x3FF; mant++ {
			s := fmt.Sprintf("%s0b%d.%010bp0", "+", lead, mant)
			m, _, err := big.ParseFloat(s, 0, 53, big.ToNearestEven)
			if err != nil {
				panic(err)
			}
			mantissa, acc := m.Float64()
			if acc != big.Exact {
				panic("not exact")
			}
			x := math.Pow(-1, float64(signbit)) * math.Pow(2, exp-bias+1) * mantissa
			bits := uint16(signbit) << 15
			bits |= uint16(exp) << 10
			bits |= uint16(mant)
			if bits == 0x8000 {
				// -zero
				d := fmt.Sprintf("{bits: 0x%04X, want: math.Copysign(0, -1)}", bits)
				ds = append(ds, d)
			} else {
				d := fmt.Sprintf("{bits: 0x%04X, want: %v}", bits, x)
				ds = append(ds, d)
			}
		}
	}
	return ds
}
