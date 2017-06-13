package main

import (
	"math"
)

func maxmin(b []byte) (max, min int) {
	max, min = 0, 255
	for i := 0; i < len(b); i++ {
		if int(b[i]) > max {
			max = int(b[i])
		}
		if int(b[i]) < min {
			min = int(b[i])
		}
	}

	return
}

func maxfreq(nrml []float64) (max, f float64) {
	max, mi := 0, 0
	off := 10
	for i, v := range nrml[off : len(nrml)/2] {
		if v > max {
			max = v
			mi = i
		}
	}

	f = float64(mi+off) * SAMP / N
	return
}

//dummy 1000hz generator
type dummy struct {
	c float64
}

func (d dummy) Read(b []byte) (n int, err error) {
	freq := 1000.0
	for i := range b {
		b[i] = byte(127*math.Sin(2*math.Pi*freq*float64(d.c)/SAMP) + 128)
		d.c += 1
	}

	return len(b), nil
}

// actually, modulus squared
func modulus(x complex64) float64 {
	return math.Sqrt(float64(real(x * complex(real(x), -imag(x)))))
}

// toss out the first few, and the top half
func maxminc(s []complex64) (max, min float64, maxi, mini int) {
	max, min = 0, math.MaxFloat64
	for i, x := range s[:len(s)/2] {
		v := modulus(x)
		if v > max {
			max, maxi = v, i
		}
		if v < min {
			min, mini = v, i
		}
	}
	return
}

func binf(i int) float64 {
	return float64(i) * SAMP / N
}

func fbin(f float64) int {
	return int(f * N / SAMP)
}

func add(xs ...float64) (total float64) {
	for _, v := range xs {
		total += v
	}
	return
}
