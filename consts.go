package main

// N == size of FFT window
// BufSize == Size of Buffer
// NBuf == number of buffers
const (
	N        = 1 << 14
	BUF_SIZE = 4 * N
	N_BUF    = 1 << 4
	CONSUME  = BUF_SIZE / N_BUF
)

const (
	RT12   = 1.0594630943592953
	BOTTOM = 55 //quite a low a
	SEMI_N = 72
)
