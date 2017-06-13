package main

import (
	"bufio"
	"github.com/mxork/lztn/overlap"
	"math"
	"os"
	"time"
	//"os/exec"
	nc "code.google.com/p/goncurses"
	"errs"
	"github.com/mxork/rbuf"
	fftw "github.com/runningwild/go-fftw/fftw32"
	"io"
)

var _ = time.Now

var xk = errs.Xk

var SAMP = 44100.0

func main() {
	// ncurses
	win, err := nc.Init()
	xk(err)
	defer nc.End()

	//nc.Cursor(0)
	//nc.Echo(false)
	//nc.CBreak(true)

	//
	bf := bufio.NewReaderSize(os.Stdin, CONSUME)

	// init various result arrays
	// two buffers to increase accuracy while decreasing latency
	rb := rbuf.New(BUF_SIZE)

	nrml := make([]float64, N)
	spec := make([]float64, SEMI_N)

	//fftw
	in := fftw.NewArray(N)
	out := fftw.NewArray(N)
	plan := fftw.NewPlan(in, out, fftw.Forward, fftw.DefaultFlag)
	defer plan.Destroy()

	//fill it up once
	_, err = io.CopyN(rb, bf, BUF_SIZE)
	xk(err)

	for cnt := 0; ; cnt++ {
		rb.Consume(CONSUME)
		_, err = io.CopyN(rb, bf, CONSUME)
		xk(err)

		// convert and transform
		// if go-fftw32 supported real-only transforms,
		// I could avoid the copy and not use a ring
		// buffer
		elems := in.Elems
		for _, buf := range rb.AsSlice() {
			fs := overlap.Bytes(buf).FloatSlice()
			for j, v := range fs {
				elems[j] = complex(v, 0)
			}
			elems = elems[len(fs):] // offset
		}

		plan.Execute()

		// normalize
		for i, v := range out.Elems {
			nrml[i] = modulus(v) / math.Sqrt(N)
		}

		// histogram calc
		max, maxf, maxfi := 0.0, 0.0, 0
		for fi, f := range midnotes {

			if fi == len(midnotes)-1 {
				break
			}

			bf, tf := f, midnotes[fi+1]
			//tighten it up a little
			Δ := tf - bf
			bf, tf = bf+Δ/4, tf-Δ/4
			bi, ti := fbin(bf), fbin(tf)

			sum := add(nrml[bi:ti]...)

			spec[fi] = sum / float64(ti-bi)
			if spec[fi] > max {
				max, maxf, maxfi = spec[fi], (tf+bf)/2, fi
			}
		}

		// display
		win.Erase()
		for i, vol := range spec {
			x := i + i/12
			win.MovePrint(0, x, noteName(i))
			win.MovePrint(1, x, "-")

			for j, a := 2, vol; a > 0 && j < 20; j, a = j+1, a-0.05 {
				win.MovePrint(j, x, ".")
			}
		}

		//supplemental info
		my, mx := win.MaxYX()
		win.MovePrint(my-10, mx-10, cnt)
		win.MovePrint(my-8, 0, maxf)
		win.MovePrint(my-7, 0, noteName(maxfi))

		if maxf > 1000 && maxf < 1300 {
			win.MovePrint(my-5, 0, "BANG")
		}

		win.NoutRefresh()

		xk(nc.Update())

		// stay up to date
		bf.Reset(os.Stdin)
	}
}
