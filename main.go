package main

import (
	"bufio"
	"fmt"
	fftw "github.com/runningwild/go-fftw/fftw32"
	"math"
	"os"
	"os/exec"
	"time"
)

var _ = fmt.Println
var _ = fftw.FFT
var _ = exec.Command
var _ = os.Stdin
var _ = time.Now

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
	for i, x := range s[20 : len(s)/2] {
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

var SAMP = 44100.0

const N = 2 << 13

func main() {
	bf := bufio.NewReaderSize(os.Stdin, int(SAMP)) // one sec
	b := make([]byte, N)
	in := fftw.NewArray(N)
	out := fftw.NewArray(N)
	nrml := make([]float64, N)
	plan := fftw.NewPlan(in, out, fftw.Forward, fftw.DefaultFlag)

	avg := 100.0
	last := time.Now()

	for _ = range time.Tick(500 * time.Millisecond) {
		n, err := bf.Read(b)
		if err != nil || n != len(b) {
			fmt.Println(err)
		}

		for i, v := range b {
			in.Elems[i] = complex((float32(v)-128)/127, 0)
		}

		plan.Execute()

		for i, v := range out.Elems {
			nrml[i] = modulus(v) / math.Sqrt(N)
		}

		sum := 0.0
		for _, v := range nrml[fbin(1100):fbin(1300)] {
			sum += v
		}
		if sum > 10*avg && time.Since(last).Seconds() > 5 {
			fmt.Println(sum, avg)
			if len(os.Args) == 2 {
				exec.Command(os.Args[1]).Run()
			}
		}

		avg = (avg + sum) / 2

		bf.Reset(os.Stdin)
	}

}
