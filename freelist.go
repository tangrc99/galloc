package galloc

import (
	"errors"
	"fmt"
	"math"
	"unsafe"
)

type spanSet map[addr]struct{}
type addr uintptr

const allocPages = 128
const allocStep = uint64(0x20000) // 128 kb

func errInvalidPointer(ptr addr) error {
	return errors.New(fmt.Sprintf("Invalid Pointer: %p", ptr))
}

//const allocStep = uint64(4096)

type freelist struct {
	ids         []addr
	freeMap     map[uint64]spanSet
	forwardMap  map[addr]uint64
	backwardMap map[addr]uint64
	pages       map[addr]Page // 从系统中分配的
	allocs      map[addr]struct{}
}

var fl *freelist

func init() {
	fl = new(freelist)
	fl.freeMap = make(map[uint64]spanSet)
	fl.forwardMap = make(map[addr]uint64)
	fl.backwardMap = make(map[addr]uint64)
	fl.pages = make(map[addr]Page)

	// startup memory pool
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
			setPageHeader(span, nt)
			f.allocs[span] = struct{}{}
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
			remain := size - uint64(nt)
			// add remain span
			f.addSpan(span+addr(nt), remain)
			println("malloc larger size")
			setPageHeader(span, nt)
			f.allocs[span] = struct{}{}
			return span + addr(pageHeaderSize)
		}
	}
	// 使用 mmap 分配内存
	npg := uint64(math.Ceil(float64(uint64(nt) / allocStep)))
	err, p := mmap(int(npg * allocStep))
	if err != nil {
		return 0
	}
	base := addr(unsafe.Pointer(&p.dataRef[0]))
	f.pages[base] = p
	f.addSpan(base+addr(nt), uint64(p.size-nt))
	println("malloc new region")
	setPageHeader(base, nt)
	f.allocs[base] = struct{}{}
	return base + addr(pageHeaderSize)
}

func (f *freelist) deallocate(ptr addr) {
	header := getPageHeader(ptr)
	start := addr(unsafe.Pointer(header))
	if _, exist := f.allocs[start]; !exist {
		panic(errInvalidPointer(ptr))
	}
	// merge existing spans
	f.mergeSpans(start, header.size)
	// try munmap pages
	if len(f.pages) > allocPages {
		// TODO: release pages
	}
}

func (f *freelist) mergeSpans(span addr, size int) {
	prev := span - 1
	next := span + addr(size)

	preSize, mergeWithPrev := f.backwardMap[prev]
	nextSize, mergeWithNext := f.forwardMap[next]
	newStart := span
	newSize := uint64(size)

	if mergeWithPrev {
		//merge with previous span
		start := prev + 1 - addr(preSize)
		f.delSpan(start, preSize)

		newStart -= addr(preSize)
		newSize += preSize
	}

	if mergeWithNext {
		// merge with next span
		f.delSpan(next, nextSize)
		newSize += nextSize
	}

	f.addSpan(newStart, newSize)
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
