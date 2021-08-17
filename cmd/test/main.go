package main

import (
	"fmt"
)

func main() {
	var a uint64 = 1 << 63
	var b uint64 = 0
	fmt.Printf("%b\n", a)
	fmt.Printf("%b\n", ^b)
	fmt.Printf("%d\n", 2+^b)
}
