// `overlap` is for mangling byte slices into different formats.
package overlap

import (
	"fmt"
	"unsafe"
)

type slice struct {
	p, l, c uintptr // ptr to data, length, capacity
}

func (s slice) String() string {
	return fmt.Sprintf("{0x%x %d %d}", s.p, s.l, s.c)
}

// Bytes is the basic type: it's methods allow an unsafe conversion to
// a slice of a different numeric type.
type Bytes []byte

// ByteSlice exists soley for completeness: the 'real' type of Bytes
// is already `[]byte`.
func (b Bytes) ByteSlice() []byte {
	return []byte(b)
}

// FloatSlice returns a slice of `float32`'s with length = len(Bytes)/4
func (b Bytes) FloatSlice() []float32 {
	// I hope this doesn't reserve any memory beyond the header
	fs := make([]float32, 0)
	fh, bh := (*slice)(unsafe.Pointer(&fs)), (*slice)(unsafe.Pointer(&b))
	fh.p, fh.l, fh.c = bh.p, bh.l/4, bh.c/4 // note, the builtin 'floor'
	return fs
}

// Int16Slice returns a slice of `int16`'s with length = len(Bytes)/2
func (b Bytes) Int16Slice() []int16 {
	fs := make([]int16, 0)
	fh, bh := (*slice)(unsafe.Pointer(&fs)), (*slice)(unsafe.Pointer(&b))
	fh.p, fh.l, fh.c = bh.p, bh.l/2, bh.c/2
	return fs
}

func (b Bytes) rawHeader() string {
	p := unsafe.Pointer(&b)
	return fmt.Sprintf("%v", *(*slice)(p))
}
