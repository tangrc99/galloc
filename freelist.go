package galloc

import "unsafe"

type spanSet map[addr]struct{}
type addr uintptr

const allocStep = uint64(0x20000)

type freelist struct {
	ids         []addr
	spans       map[addr]struct{} // 当前可用的 span
	freeMap     map[uint64]spanSet
	forwardMap  map[addr]uint64
	backwardMap map[addr]uint64

	pages map[addr]Page // 从系统中分配的
}

var fl *freelist

func init() {
	fl = new(freelist)
	fl.spans = make(map[addr]struct{})
	fl.freeMap = make(map[uint64]spanSet)
	fl.forwardMap = make(map[addr]uint64)
	fl.backwardMap = make(map[addr]uint64)
	fl.pages = make(map[addr]Page)
}

func (f *freelist) allocate(n int) addr {

	nt := n + pageHeaderSize

	// 先寻找是否有对应大小的 page
	if spans, ok := f.freeMap[uint64(nt)]; ok {
		for span := range spans {
			// 删除对应 page 的记录
			delete(spans, span)
			// TODO: 记录本次分配
			println("malloc proper size")
			setPageHeader(span, n)
			return span + addr(pageHeaderSize)
		}
	}

	// 再寻找稍大的 page
	for size, spans := range f.freeMap {
		if size < uint64(nt) {
			continue
		}
		for span := range spans {
			// remove the initial
			f.delSpan(span, size)

			//f.allocs[span] = txid
			remain := size - uint64(nt)

			// add remain span
			f.addSpan(span+addr(nt), remain)

			for i := addr(0); i < addr(nt); i++ {
				delete(f.spans, span+i)
			}
			println("malloc larger size")
			setPageHeader(span, n)

			return span + addr(pageHeaderSize)
		}
	}

	// 使用 mmap 分配内存
	npg := uint64(nt)/allocStep + 1
	err, p := mmap(int(npg * allocStep))
	if err != nil {
		return 0
	}
	base := addr(unsafe.Pointer(&p.dataRef[0]))
	f.pages[base] = p
	f.addSpan(base+addr(nt), uint64(p.size-nt))
	println("malloc new region")
	setPageHeader(base, n)
	return base + addr(pageHeaderSize)
}

func (f *freelist) deallocate(ptr addr, n int) {

}

func (f *freelist) mergeSpans() {

}

func (f *freelist) addSpan(start addr, size uint64) {
	f.backwardMap[start-1+addr(size)] = size
	f.forwardMap[start] = size
	if _, ok := f.freeMap[size]; !ok {
		f.freeMap[size] = make(map[addr]struct{})
	}

	f.freeMap[size][start] = struct{}{}
}

func (f *freelist) delSpan(start addr, size uint64) {
	delete(f.forwardMap, start)
	delete(f.backwardMap, start+addr(size-1))
	delete(f.freeMap[size], start)
	if len(f.freeMap[size]) == 0 {
		delete(f.freeMap, size)
	}
}
