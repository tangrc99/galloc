package galloc

import (
	"errors"
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"unsafe"
)

type spanSet map[addr]struct{}
type addr uintptr

const (
	// memoryPageSize 内存页大小
	memoryPageSize = 0x1000 //4kb
	// allocStep 是最小分配单元
	allocStep = uint64(32 * memoryPageSize) // 128kb
	// maxAllocPages 是常驻的最大页数
	maxAllocPages = 128
	// startupPages 是初始化时分配的页数
	startupPages = 1
)

func errInvalidPointer(ptr addr) error {
	return errors.New(fmt.Sprintf("Invalid Pointer: %x", ptr))
}

type freelist struct {
	nShard      int
	freeMap     map[uint64]spanSet // 长度 - 地址集合
	forwardMap  map[addr]uint64    // 正向查找
	backwardMap map[addr]uint64    // 反向查找
	pages       map[addr]Page      // 从系统中分配的
	allocMap    map[addr]struct{}  // 分配给用户的内存
	lk          spinLock
}

func (f *freelist) allocate(n int) addr {
	// TODO: 修改内存对齐策略
	nt := pow2upper(n + pageHeaderSize)
	f.lk.Lock()
	// 先寻找是否有对应大小的 page
	if spans, ok := f.freeMap[uint64(nt)]; ok {
		for span := range spans {
			// 删除对应 page 的记录
			f.delSpan(span, uint64(nt))
			f.allocMap[span] = struct{}{}
			f.lk.Unlock()
			setPageHeader(span, nt, f.nShard)
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
			f.allocMap[span] = struct{}{}
			f.lk.Unlock()
			setPageHeader(span, nt, f.nShard)
			return span + addr(pageHeaderSize)
		}
	}
	// 使用 mmap 分配内存
	npg := uint64(math.Ceil(float64(nt) / float64(allocStep)))
	err, p := mmap(int(npg * allocStep))
	if err != nil {
		return 0
	}
	base := addr(unsafe.Pointer(&p.dataRef[0]))
	f.pages[base] = p
	f.addSpan(base+addr(nt), uint64(p.size-nt))
	f.allocMap[base] = struct{}{}
	f.lk.Unlock()
	setPageHeader(base, nt, f.nShard)
	return base + addr(pageHeaderSize)
}

func (f *freelist) deallocate(ptr addr) {
	header := getPageHeader(ptr)
	start := addr(unsafe.Pointer(header))
	f.lk.Lock()
	if _, exist := f.allocMap[start]; !exist {
		panic(errInvalidPointer(ptr))
	}
	delete(f.allocMap, start)
	// merge existing spans
	f.mergeSpans(start, header.size)
	// release pages if needed
	if len(f.pages) > maxAllocPages {
		for span, pg := range f.pages {
			if sz := f.forwardMap[span]; sz >= uint64(pg.size) {
				f.delSpan(span, sz)
				if sz > allocStep {
					f.addSpan(span+addr(allocStep), sz-allocStep)
				}
				delete(f.pages, span)
				f.lk.Unlock()
				//println("munmap")
				if err := munmap(pg.dataRef); err != nil {
					panic(err)
				}
				return
			}
		}
	}
	f.lk.Unlock()
}

func (f *freelist) mergeSpans(span addr, size int) {
	prev := span - 1
	next := span + addr(size)

	preSize, mergeWithPrev := f.backwardMap[prev]
	nextSize, mergeWithNext := f.forwardMap[next]
	newStart := span
	newSize := uint64(size)
	//if mergeWithPrev {
	if _, exist := f.pages[span]; mergeWithPrev && !exist {
		//merge with previous span, when start is not a page
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

type allocator []*freelist

var alloc allocator

func (a *allocator) allocate(n int) addr {
	nShard := rand.Intn(len(alloc))
	fl := alloc[nShard]
	return fl.allocate(n)
}

func (a *allocator) deallocate(ptr addr) {
	nShard := getPageHeader(ptr).nShard
	fl := alloc[nShard]
	fl.deallocate(ptr)
}

func init() {
	alloc = make([]*freelist, runtime.NumCPU())

	for i := range alloc {
		alloc[i] = new(freelist)
		alloc[i].freeMap = make(map[uint64]spanSet)
		alloc[i].forwardMap = make(map[addr]uint64)
		alloc[i].backwardMap = make(map[addr]uint64)
		alloc[i].pages = make(map[addr]Page)
		alloc[i].allocMap = make(map[addr]struct{})
		alloc[i].nShard = i
		// startup memory pool
		for j := 0; j < startupPages; j++ {
			err, p := mmap(int(allocStep))
			if err != nil {
				panic(err)
			}
			base := addr(unsafe.Pointer(&p.dataRef[0]))
			alloc[i].addSpan(base, allocStep)
		}
	}
}
