package main

import (
	"fmt"
	"github.com/liyue201/gostl/algorithm"
	"github.com/liyue201/gostl/ds/deque"
)

func main() {
	a := deque.New()
	for i := 0; i < 9; i++ {
		a.PushBack(i)
	}
	fmt.Printf("%v\n", a)

	algorithm.Swap(a.First(), a.Last())
	fmt.Printf("%v\n", a)

	algorithm.Reverse(a.Begin(), a.End())
	fmt.Printf("%v\n", a)
}
