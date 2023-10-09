package example

import (
	"container/list"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
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

type smallStruct struct {
	aa [1024]byte
}

type bigStruct struct {
	bb [65536]byte
}

func getBigList(nNode int, value any) *List {
	l := NewList()
	for i := 0; i < nNode; i++ {
		l.PushBack(value)
		if i%10000 == 0 {
			fmt.Printf("allocated nodes %d\n", i)
		}
	}
	return l
}

func getBigList2(nNode int, value any) *list.List {
	l := list.New()
	for i := 0; i < nNode; i++ {
		l.PushBack(value)
		if i%10000 == 0 {
			fmt.Printf("allocated nodes %d\n", i)
		}
	}
	return l
}

func TestListSmallElement(t *testing.T) {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	l := make([]*List, 10)
	for i := range l {
		l[i] = getBigList(100000, smallStruct{})
		if i%1000 == 0 {
			fmt.Printf("allocated %d\n", i)
		}
	}
	s := time.Now()
	runtime.GC()
	fmt.Printf("list1 gc time: %dus\n", time.Since(s).Microseconds())

}

func TestListSmallElement2(t *testing.T) {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	l := make([]*list.List, 10)
	for i := range l {
		l[i] = getBigList2(100000, smallStruct{})
		if i%1000 == 0 {
			fmt.Printf("allocated %d\n", i)
		}
	}
	s := time.Now()
	runtime.GC()
	fmt.Printf("list1 gc time: %dus\n", time.Since(s).Microseconds())

}

func TestListBigElement(t *testing.T) {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()

	l := make([]*List, 10)
	for i := range l {
		l[i] = getBigList(100000, bigStruct{})
		if i%1000 == 0 {
			fmt.Printf("allocated %d\n", i)
		}
	}

	s := time.Now()
	runtime.GC()
	fmt.Printf("list1 gc time: %dus\n", time.Since(s).Microseconds())

}

func TestListBigElement2(t *testing.T) {
	go func() {
		log.Println(http.ListenAndServe(":6060", nil))
	}()
	l := make([]*list.List, 10)
	for i := range l {
		l[i] = getBigList2(100000, bigStruct{})
		if i%1000 == 0 {
			fmt.Printf("allocated %d\n", i)
		}
	}

	s := time.Now()
	runtime.GC()
	fmt.Printf("list2 gc time: %dus\n", time.Since(s).Microseconds())
}
