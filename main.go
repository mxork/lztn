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

const SAMP = 44100
const SMPI = 1 / 44100
const N = 2 << 14

func main() {
	bf := bufio.NewReaderSize(os.Stdin, SAMP) // one sec
	b := make([]byte, N)
	in := fftw.NewArray(N)
	out := fftw.NewArray(N)
	plan := fftw.NewPlan(in, out, fftw.Forward, fftw.DefaultFlag)

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
			fmt.Println(float64(i)*SAMP/N, real(v*complex(real(v), -imag(v))))
		}

		bf.Reset(os.Stdin)
	}

}
