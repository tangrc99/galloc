package galloc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/tmthrgd/go-memset"
	"runtime"
	"testing"
	"time"
	"unsafe"
)

func TestGalloc(t *testing.T) {
	type A struct {
		Val int
	}
	aa := New[A]()
	aa.Val = 1
	Delete(aa)
}

func TestAllocate(t *testing.T) {
	type A struct {
		Val int
	}
	ptr := alloc.allocate(1024)
	aptr := (*A)(unsafe.Pointer(ptr))
	aptr.Val = 1
	alloc.deallocate(ptr)

	ptr = alloc.allocate(1024)
	alloc.deallocate(ptr)
}

func TestFree(t *testing.T) {
	a1 := alloc.allocate(int(float64(allocStep)*1.5) - pageHeaderSize)
	a2 := alloc.allocate(int(float64(allocStep)*0.5) - pageHeaderSize)
	alloc.deallocate(a1)
	alloc.deallocate(a2)
}

func TestFree2(t *testing.T) {
	a1 := alloc.allocate(int(float64(allocStep)*0.5) - pageHeaderSize)
	a2 := alloc.allocate(int(float64(allocStep)*0.5) - pageHeaderSize)
	a3 := alloc.allocate(int(float64(allocStep)*0.5) - pageHeaderSize)
	a4 := alloc.allocate(int(float64(allocStep)*0.5) - pageHeaderSize)
	alloc.deallocate(a2)
	alloc.deallocate(a3)
	alloc.deallocate(a4)
	alloc.deallocate(a1)
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

func BenchmarkMalloc(b *testing.B) {
	s := Malloc(1024)
	Free(s)

	nTest := 10000
	b.ResetTimer()
	for i := 0; i < nTest; i++ {
		s := Malloc(1024)
		Free(s)
	}
}

func TestNew(t *testing.T) {

	type A struct {
		a int
	}
	l := make([]*A, 100000)

	for i := range l {
		l[i] = New[A]()
	}
	s := time.Now()
	runtime.GC()
	fmt.Printf("galloc.New() gc time: %dus\n", time.Since(s).Microseconds())
}

func TestNew2(t *testing.T) {
	type A struct {
		a int
	}
	l := make([]*A, 100000)

	for i := range l {
		l[i] = new(A)
	}

	s := time.Now()
	runtime.GC()
	fmt.Printf("new() gc time: %dus\n", time.Since(s).Microseconds())
}
