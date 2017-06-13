package overlap

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestBasic(t *testing.T) {
	b := Bytes(make([]byte, 8, 12))
	fs := b.FloatSlice()

	fmt.Println(b.rawHeader())
	fmt.Println(unsafe.Sizeof(b))
	fmt.Println(fs)
}
