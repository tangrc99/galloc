package galloc

import (
	"fmt"
	"testing"
	"unsafe"
)

func Testgalloc(t *testing.T) {
	aa := New[A]()
	aa.Val = 1
	println(aa)

	println(aa.Val)
	//Free[A](aa)

	base := uintptr(unsafe.Pointer(aa)) - uintptr(pageHeaderSize)
	bptr := (*pageHeader)(unsafe.Pointer(base))
	fmt.Printf("%x\n", base)
	println(bptr.size)
}

func TestAllocate(t *testing.T) {
	ptr := fl.allocate(1024)
	aptr := (*A)(unsafe.Pointer(ptr))
	aptr.Val = 1
	println(aptr.Val)
	ptr = fl.allocate(1024)
}
