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

## 压力测试

``` bash
Usage of est:
  -H string
        http headers
  -P uint
        LPS 默认为1 (default 1)
  -T int
        TimeoutNS 以 MS 计 (default -1)
  -U string
        http API 地址(仅支持http) 必须添加 http://
  -X string
        http 请求方式 (default "GET")
  -d string
        http body
  -f string
        http body from file 输入文件路径 此选项会覆盖 -d
  -t int
        测试持续时间 以 S 计时 (default -1)
```

``` bash
go build && go install

est -U http://127.0.0.1:9426/api/v1/login -T 50 -t 5 -X POST -P 1000 -d '{"username":"yixingwei","password":"123456"}'

est -U http://127.0.0.1:9426/api/v1/login -T 50 -t 5 -X POST -P 1000 -f ./post.json
```

