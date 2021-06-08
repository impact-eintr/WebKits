# WebKits
常用工具 以及 数据结构

## 令牌桶

```go
package etbf

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

const (
	SIZE = 5
)

func TestFetchtoken(t *testing.T) {

	wp := sync.WaitGroup{}

	tbf1 := Newtbf(time.Second, 5, 100)
	if tbf1 == nil {
		log.Fatalln(errors.New("无法初始化"))
	}

	file1, err := os.Open("./log1")
	if err != nil {
		log.Println(err)
	}
	defer file1.Close()

	wp.Add(2)

	go func(f io.ReadWriter) {
		for {
			size, err := tbf1.Fetchtoken(SIZE)
			if err != nil {
				log.Fatalln(err)
			}

			buf := make([]byte, size)
			n, err := f.Read(buf)
			if n == 0 && err == io.EOF {
				tbf1.Destory()
				break
			}

			fmt.Print("1" + string(buf))
		}

		wp.Done()
	}(file1)

	tbf2 := Newtbf(2*time.Second, 10, 100)
	if tbf2 == nil {
		log.Fatalln(errors.New("无法初始化"))
	}
	file2, err := os.Open("./log2")
	if err != nil {
		log.Println(err)
	}
	defer file2.Close()

	go func(f io.ReadWriter) {
		for {
			size, err := tbf2.Fetchtoken(SIZE)
			if err != nil {
				log.Fatalln(err)
			}

			buf := make([]byte, size)
			n, err := f.Read(buf)
			if n == 0 && err == io.EOF {
				tbf2.Destory()
				break
			}

			fmt.Print("2" + string(buf))
		}

		wp.Done()
	}(file2)

	wp.Wait()

}
```

