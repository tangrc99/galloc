# galloc

`galloc`是基于 bbolt 中内存页哈希分配策略实现的一个玩具 Demo。为了规避 Golang 中的复杂结构体 GC 问题，galloc 实现了一个简单的内存池，内存池中的分配器使用`mmap`实现，内存页的分配策略借鉴了`bbolt`中的分配策略。由`galloc`分配的内存不参与 Golang 的 GC，可以用于构建链表、树等指针较多的数据结构。

## Usage

```
// 使用 uintptr 构建的链表节点
type listNode struct {
    val any
    next uintptr
}

// 分配链表节点
node1 := galloc.New[listNode]()
node2 := galloc.New[listNode]()
node1.next = uintptr(unsafe.Pointer(node2))

// 析构对象
galloc.Free(node2)
galloc.Free(node1)
```

## Note

以指针形式保留从`galloc`中获取的内存并不能够规避 Golang 的 GC，因此仅仅用该接口替代`new()`接口并不能改善 GC。只有使用`uintptr`这一数据形式来保留`galloc`中分配出的内存地址才能够起到减少 GC 的效果。

## Performance

简单测试了链表形式下`galloc`的分配速度以及 GC 改善效果，baseline 为 Golang 内存分配器。在完成 10 个 100000 节点链表的分配后，立刻进入 GC：

|                   | 分配速度/s | GC 时长/us |
| ----------------- | ---------- | ---------- |
| galloc            | 14.07      | 269        |
| Golang 内存分配器 | 60.23      | 57479      |

基于`galloc`分配内存构建的链表隐藏了内部的指针，所以在 GC 时只需要寻址扫描链表的头结点，而使用官方的内存分配器则需要对内部节点完整扫描，因此 GC 时间较长。