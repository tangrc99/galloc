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

var lk spinLock

func BenchmarkMallocWithSpinLock(b *testing.B) {
	s := Malloc(1024)
	Free(s)

	nTest := 10000
	b.ResetTimer()
	for i := 0; i < nTest; i++ {
		lk.Lock()
		s := Malloc(1024)
		lk.Unlock()
		lk.Lock()
		Free(s)
		lk.Unlock()
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
