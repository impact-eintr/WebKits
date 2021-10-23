package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c)
	fmt.Println("start..")
	s := <-c

	switch s {
	case syscall.SIGTERM:
		fmt.Println("End...", s)
	}
}
