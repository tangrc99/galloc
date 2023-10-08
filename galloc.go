package galloc

import (
	"github.com/tmthrgd/go-memset"
	"unsafe"
)

func New[T any]() *T {
	var t *T
	n := unsafe.Sizeof(*t)
	t = (*T)(unsafe.Pointer(fl.allocate(int(n))))
	memset.Memset((*[maxMapSize]byte)(unsafe.Pointer(t))[0:n], 0)
	return t
}

func Delete[T any](obj *T) {
	fl.deallocate(addr(unsafe.Pointer(obj)))
}

func Malloc(n int) []byte {
	ptr := fl.allocate(n)
	return (*[maxMapSize]byte)(unsafe.Pointer(ptr))[0:n]
}

func Free(bytes []byte) {
	fl.deallocate(addr(unsafe.Pointer(&bytes[0])))
}

func ToAddr(ptr *any) uintptr {
	return uintptr(unsafe.Pointer(ptr))
}

func ToPtr[T any](ptr uintptr) *T {
	return (*T)(unsafe.Pointer(ptr))
}
