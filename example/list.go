package example

import (
	"galloc"
	"unsafe"
)

type listNode struct {
	Value any
	next  uintptr
	prev  uintptr
}

func newListNode(val any) *listNode {
	ptr := galloc.New[listNode]()
	ptr.Value = val
	ptr.prev = 0
	ptr.next = 0
	return ptr
}

func (n *listNode) nextNode() *listNode {
	return (*listNode)(unsafe.Pointer(n.next))
}

func (n *listNode) prevNode() *listNode {
	return (*listNode)(unsafe.Pointer(n.prev))
}

type List struct {
	sentinel *listNode
	nNode    int
}

func NewList() *List {
	sentinel := newListNode(0)
	sentinel.next = uintptr(unsafe.Pointer(sentinel))
	sentinel.prev = uintptr(unsafe.Pointer(sentinel))
	return &List{
		sentinel: sentinel,
		nNode:    0,
	}
}

func (l *List) senti() uintptr {
	return uintptr(unsafe.Pointer(l.sentinel))
}

func (l *List) Empty() bool {
	return l.sentinel.nextNode() == l.sentinel
}

func (l *List) Front() any {
	if l.Empty() {
		return nil
	}
	return l.sentinel.nextNode().Value
}

func (l *List) Back() any {
	if l.Empty() {
		return nil
	}
	return l.sentinel.prevNode().Value
}

func (l *List) PushBack(val any) {
	node := newListNode(val)
	old := l.sentinel.prevNode()
	old.next = uintptr(unsafe.Pointer(node))
	l.sentinel.prev = uintptr(unsafe.Pointer(node))
	node.prev = uintptr(unsafe.Pointer(old))
	node.next = uintptr(unsafe.Pointer(l.sentinel))
	l.nNode += 1
}

func (l *List) PushFront(val any) {
	node := newListNode(val)
	old := l.sentinel.nextNode()
	old.prev = uintptr(unsafe.Pointer(node))
	l.sentinel.next = uintptr(unsafe.Pointer(node))
	node.next = uintptr(unsafe.Pointer(old))
	node.prev = uintptr(unsafe.Pointer(l.sentinel))
	l.nNode += 1
}

func (l *List) PopBack() any {
	if l.Empty() {
		return nil
	}
	pop := l.sentinel.prevNode()
	tail := pop.prevNode()
	tail.next = uintptr(unsafe.Pointer(l.sentinel))
	l.sentinel.prev = uintptr(unsafe.Pointer(tail))
	l.nNode -= 1
	v := pop.Value
	galloc.Delete(pop)
	return v
}

func (l *List) PopFront() any {
	if l.Empty() {
		return nil
	}
	pop := l.sentinel.nextNode()
	front := pop.nextNode()
	front.prev = uintptr(unsafe.Pointer(l.sentinel))
	l.sentinel.next = uintptr(unsafe.Pointer(front))
	l.nNode -= 1
	v := pop.Value
	galloc.Delete(pop)
	return v
}

func (l *List) Size() int {
	return l.nNode
}
