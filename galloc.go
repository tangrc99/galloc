package galloc

import (
	"github.com/tmthrgd/go-memset"
	"unsafe"
)

// New creates an object and return its ptr. The object must be freed using Delete
func New[T any]() *T {
	var t *T
	n := unsafe.Sizeof(*t)
	t = (*T)(unsafe.Pointer(alloc.allocate(int(n))))
	memset.Memset((*[maxMapSize]byte)(unsafe.Pointer(t))[0:n], 0)
	return t
}

// Delete frees an object created by New. If obj is invalid, a panic will be thrown
func Delete[T any](obj *T) {
	alloc.deallocate(addr(unsafe.Pointer(obj)))
}

// Malloc allocates n-bytes memory segment, which must be freed using Free
func Malloc(n int) []byte {
	ptr := alloc.allocate(n)
	return (*[maxMapSize]byte)(unsafe.Pointer(ptr))[0:n]
}

// Free frees the memory segment allocated by Malloc
func Free(bytes []byte) {
	alloc.deallocate(addr(unsafe.Pointer(&bytes[0])))
}

// ToAddr converts a *T object allocated by galloc to uintptr
func ToAddr(ptr *any) uintptr {
	return uintptr(unsafe.Pointer(ptr))
}

// ToPtr converts a uintptr object to *T
func ToPtr[T any](ptr uintptr) *T {
	return (*T)(unsafe.Pointer(ptr))
}
