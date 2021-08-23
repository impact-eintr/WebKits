package main

import "fmt"

func main() {
	m := make(map[int]int, 1<<10)
	for i := 0; i < 1024; i++ {
		m[i] = i << 1
	}

	for k := range m {
		go func(key int) {
			fmt.Println(key)
		}(k)
	}
}
