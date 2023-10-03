package galloc

import "unsafe"

const maxMapSize = 0xFFFFFFFF
const pageHeaderSize = int(unsafe.Sizeof(pageHeader{}))

type Page struct {
	dataRef []byte
	data    *[maxMapSize]byte
	size    int
}

type pageHeader struct {
	size int // 析构的时候用来确认大小
}

func setPageHeader(ptr addr, n int) {
	(*pageHeader)(unsafe.Pointer(ptr)).size = n
}
