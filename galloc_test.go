package galloc

import (
	"testing"
	"unsafe"
)

type A struct {
	Val int
}

func TestGalloc(t *testing.T) {
	aa := New[A]()
	aa.Val = 1
	println(aa)

	println(aa.Val)
	//Delete[A](aa)

	//base := uintptr(unsafe.Pointer(aa)) - uintptr(pageHeaderSize)
	//bptr := (*pageHeader)(unsafe.Pointer(base))
	//fmt.Printf("%x\n", base)
	//println(bptr.size)

	Delete(aa)
}

func TestAllocate(t *testing.T) {
	ptr := fl.allocate(1024)
	aptr := (*A)(unsafe.Pointer(ptr))
	aptr.Val = 1
	println(aptr.Val)
	ptr = fl.allocate(1024)
}

func TestFree(t *testing.T) {
	a1 := fl.allocate(4096 - 8)
	a2 := fl.allocate(4096 - 8)
	a3 := fl.allocate(4096 - 8)
	a4 := fl.allocate(4096 - 8)
	fl.deallocate(a1)
	fl.deallocate(a2)
	fl.deallocate(a3)
	fl.deallocate(a4)

}
