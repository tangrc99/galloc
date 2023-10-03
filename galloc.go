package galloc

import (
	"unsafe"
)

func New[T any]() *T {
	var t *T
	n := unsafe.Sizeof(*t)
	// TODO: 这里要对齐为 pow2
	t = (*T)(unsafe.Pointer(fl.allocate(int(n))))
	return t
}

func Delete[T any](obj *T) {
	fl.deallocate(addr(unsafe.Pointer(obj)))
}

func Malloc(n int) []byte {
	ptr := fl.allocate(n)
	return *(*[]byte)(unsafe.Pointer(ptr))
}

func Free(bytes []byte) {
	fl.deallocate(addr(unsafe.Pointer(&bytes[0])))
}
