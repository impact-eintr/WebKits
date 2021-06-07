package etbf

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"testing"
)

const (
	SIZE = 3
)

func TestFetchtoken(t *testing.T) {

	wp := sync.WaitGroup{}

	tbf := Newtbf(10, 100)
	if tbf == nil {
		log.Fatalln(errors.New("无法初始化"))
	}

	file, err := os.Open("./log")
	if err != nil {
		log.Println(err)
	}

	wp.Add(1)

	go func(f io.ReadWriter) {
		for {
			size, err := tbf.Fetchtoken(SIZE)
			if err != nil {
				log.Fatalln(err)
			}

			buf := make([]byte, size)
			n, err := f.Read(buf)
			if n == 0 && err == io.EOF {
				break
			}

			fmt.Print(string(buf))
		}

		wp.Done()
	}(file)

	wp.Wait()

}
