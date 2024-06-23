package mq

import (
	"sync"
)

type LinkedQueue[T any] struct {
	sync.Mutex
	head     *Node[T]
	last     *Node[T]
	capacity int
}

type Node[T any] struct {
	item T
	next *Node[T]
}

func NewNode[T any](item T, next *Node[T]) *Node[T] {
	return &Node[T]{item: item, next: next}
}

func NewLinkedQueue[T any]() *LinkedQueue[T] {
	return &LinkedQueue[T]{
		head:     &Node[T]{},
		last:     &Node[T]{},
		capacity: 0,
	}
}

func (l *LinkedQueue[T]) Push(item T) {
	l.Lock()
	defer l.Unlock()
	newNode := NewNode[T](item, nil)
	if l.head.next == nil {
		l.head.next = newNode
		l.last.next = newNode
	} else {
		l.last.next.next = newNode
		l.last.next = newNode
	}
	l.capacity++
}

func (l *LinkedQueue[T]) Pop() T {
	if l.capacity <= 0 {
		var empty T
		return empty
	} else {
		l.Lock()
		defer l.Unlock()
		t := l.head.next
		l.head.next = t.next
		if l.last.next == t {
			l.last.next = nil
		}
		l.capacity--
		return t.item
	}
}

func (l *LinkedQueue[T]) GetCapacity() int {
	return l.capacity
}
