package proxy

import (
	"fmt"
	"testing"
	"unsafe"
)

func TestUnsafePoint(t *testing.T) {
	a := []byte{
		1, 2,
	}
	integer := *(*uint16)(unsafe.Pointer(&a[0]))
	fmt.Println(integer)
}
