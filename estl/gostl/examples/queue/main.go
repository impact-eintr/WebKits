package main

import (
	"fmt"
	"github.com/liyue201/gostl/ds/queue"
	"sync"
	"time"
)

func example1() {
	fmt.Printf("example1:\n")
	q := queue.New()
	for i := 0; i < 5; i++ {
		q.Push(i)
	}
	for !q.Empty() {
		fmt.Printf("%v\n", q.Pop())
	}
}

// using list as container
func example2() {
	fmt.Printf("example2:\n")
	q := queue.New(queue.WithListContainer())
	for i := 0; i < 5; i++ {
		q.Push(i)
	}
	for !q.Empty() {
		fmt.Printf("%v\n", q.Pop())
	}
}

// goroutine-save
func example3() {
	fmt.Printf("example3:\n")

	s := queue.New(queue.WithGoroutineSafe())
	sw := sync.WaitGroup{}
	sw.Add(2)
	go func() {
		defer sw.Done()
		for i := 0; i < 10; i++ {
			s.Push(i)
			time.Sleep(time.Microsecond * 100)
		}
	}()

	go func() {
		defer sw.Done()
		for i := 0; i < 10; {
			if !s.Empty() {
				val := s.Pop()
				fmt.Printf("%v\n", val)
				i++
			} else {
				time.Sleep(time.Microsecond * 100)
			}
		}
	}()
	sw.Wait()
}

func main() {
	example1()
	example2()
	example3()
}
