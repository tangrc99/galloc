package example

import (
	"container/list"
	"fmt"
	"github.com/stretchr/testify/assert"
	"runtime"
	"testing"
	"time"
)

func TestList(t *testing.T) {

	l := NewList()
	l.PushBack(1)
	l.PushBack(2)
	assert.Equal(t, 1, l.Front())
	assert.Equal(t, 2, l.Back())

	l.PopFront()
	l.PopFront()
}

func getBigList(nNode int) *List {
	l := NewList()
	for i := 0; i < nNode; i++ {
		l.PushBack(1)
	}
	return l
}

func getBigList2(nNode int) *list.List {
	l := list.New()
	for i := 0; i < nNode; i++ {
		l.PushBack(1)
	}
	return l
}

func TestListGC(t *testing.T) {
	l := make([]*List, 1000)
	for i := range l {
		l[i] = getBigList(10000)
		if i%1000 == 0 {
			fmt.Printf("allocated %d\n", i)
		}
	}
	s := time.Now()
	runtime.GC()
	fmt.Printf("list1 gc time: %dus\n", time.Since(s).Microseconds())
}

func TestListGC2(t *testing.T) {
	l := make([]*list.List, 1000)
	for i := range l {
		l[i] = getBigList2(10000)
		if i%1000 == 0 {
			fmt.Printf("allocated %d\n", i)
		}
	}
	s := time.Now()
	runtime.GC()
	fmt.Printf("list2 gc time: %dus\n", time.Since(s).Microseconds())
}
