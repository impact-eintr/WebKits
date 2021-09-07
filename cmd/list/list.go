package main

import (
	"fmt"
)

type Intrusive interface {
	Next() Intrusive
	Prev() Intrusive
	AddNext(Intrusive)
	AddPrev(Intrusive)
}

type List struct {
	prev Intrusive
	next Intrusive
}

func (l *List) Next() Intrusive {
	return l.next
}

func (l *List) Prev() Intrusive {
	return l.prev
}

func (l *List) AddNext(i Intrusive) {
	l.next = i
}

func (l *List) AddPrev(i Intrusive) {
	l.prev = i
}

func (l *List) Front() Intrusive {
	return l.prev
}

func (l *List) Back() Intrusive {
	return l.next
}

func (l *List) PushFront(e Intrusive) {
	e.AddPrev(nil)
	e.AddNext(l.prev) // 指向第一个节点

	if l.prev != nil {
		l.prev.AddPrev(e)
	} else {
		l.next = e
	}
	l.prev = e

}

func (l *List) PushBack(e Intrusive) {
	e.AddPrev(l.next) // 指向最后一个节点
	e.AddNext(nil)

	if l.next != nil {
		l.next.AddNext(e)
	} else {
		l.prev = e
	}

	l.next = e

}

// InsertAfter inserts e after b.
func (l *List) InsertAfter(e, b Intrusive) {
	a := b.Next()
	e.AddNext(a)
	e.AddPrev(b)
	b.AddNext(e)

	if a != nil {
		a.AddPrev(e)
	} else {
		l.next = e
	}
}

// InsertBefore inserts e before a.
func (l *List) InsertBefore(e, a Intrusive) {
	b := a.Prev()
	e.AddNext(a)
	e.AddPrev(b)
	a.AddPrev(e)

	if b != nil {
		b.AddNext(e)
	} else {
		l.prev = e
	}
}

// Remove removes e from l.
func (l *List) Remove(e Intrusive) {
	prev := e.Prev()
	next := e.Next()

	if prev != nil {
		prev.AddNext(next)
	} else {
		l.prev = next
	}

	if next != nil {
		next.AddPrev(prev)
	} else {
		l.next = prev
	}
}

func main() {
	type E struct {
		List
		data int
	}
	// Create a new list and put some numbers in it.
	l := List{}
	e4 := &E{data: 4}
	e3 := &E{data: 3}
	e2 := &E{data: 2}
	e1 := &E{data: 1}

	l.PushFront(e4)
	l.PushFront(e3)
	l.PushFront(e2)
	l.PushFront(e1)

	fmt.Println()

	for e := l.Front(); e != nil; e = e.Next() {
		fmt.Printf("e: %+v\n", e)
		fmt.Printf("data: %d\n", e.(*E).data)
	}

}
