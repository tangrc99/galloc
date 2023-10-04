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
	size   int // 析构的时候用来确认大小
	nShard int // shard 序号
}

func setPageHeader(ptr addr, n int) {
	(*pageHeader)(unsafe.Pointer(ptr)).size = n
}

func getPageHeader(ptr addr) *pageHeader {
	return (*pageHeader)(unsafe.Pointer(ptr - addr(pageHeaderSize)))
}
