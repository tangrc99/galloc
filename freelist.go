package galloc

type pgSet map[uint64]struct{}
type pgid uint64

type freelist struct {
	ids         []pgid
	pages       map[pgid]Page
	freeMap     map[uint64]pgSet // key is the size of continuous pages(span), value is a set which contains the starting pgids of same size
	forwardMap  map[pgid]uint64  // key is start pgid, value is its span size
	backwardMap map[pgid]uint64  // key is end pgid, value is its span size
}

func (f *freelist) allocate(size int) []byte {

	// 先寻找是否有对应大小的 page
	if pg, ok := f.freeMap[uint64(size)]; ok {
		for pid := range pg {
			// 删除对应 page 的记录

			// 记录本次分配
			return pid
		}
	}

	// 再寻找稍大的 page
}

func (f *freelist) mergeSpans() {

}

func (f *freelist) deallocate() {

}
