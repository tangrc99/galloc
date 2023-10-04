package galloc

import (
	"github.com/stretchr/testify/assert"
	"github.com/tmthrgd/go-memset"
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
	fl.deallocate(ptr)

	ptr = fl.allocate(1024)
	fl.deallocate(ptr)

}

func TestFree(t *testing.T) {
	a1 := fl.allocate(int(float64(allocStep)*1.5) - pageHeaderSize)
	a2 := fl.allocate(int(float64(allocStep)*0.5) - pageHeaderSize)

	fl.deallocate(a1)
	fl.deallocate(a2)

}

func TestFree2(t *testing.T) {
	a1 := fl.allocate(int(float64(allocStep)*0.5) - pageHeaderSize)
	a2 := fl.allocate(int(float64(allocStep)*0.5) - pageHeaderSize)
	a3 := fl.allocate(int(float64(allocStep)*0.5) - pageHeaderSize)
	a4 := fl.allocate(int(float64(allocStep)*0.5) - pageHeaderSize)
	fl.deallocate(a2)
	fl.deallocate(a3)
	fl.deallocate(a4)
	fl.deallocate(a1)
}

func TestBZero(t *testing.T) {
	b := Malloc(1024)
	assert.Equal(t, 1024, len(b))
	memset.Memset(b, 0)
	for i := range b {
		assert.Equal(t, uint8(0), b[i])
	}
	Free(b)
}
