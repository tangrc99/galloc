package galloc

import (
	"unsafe"
)

type A struct {
	Val int
}

//func upper2pow() int {
//
//}

func New[T any]() *T {
	var t *T
	n := unsafe.Sizeof(*t)
	// TODO: 这里要对齐为 pow2
	t = (*T)(unsafe.Pointer(fl.allocate(int(n))))
	return t
}

func Free[T any](obj *T) {
}
