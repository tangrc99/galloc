package galloc

import "unsafe"

const maxMapSize = 0xFFFFFFFF
const pageHeaderSize = 10

type Page struct {
	id      uint64
	dataRef []byte
	data    *[maxMapSize]byte
	size    int
}

type pageHeader struct {
	size int // 析构的时候用来确认大小

}

func getPageHeader(pg *Page) *pageHeader {
	return (*pageHeader)(unsafe.Pointer(&pg.dataRef[0]))
}
