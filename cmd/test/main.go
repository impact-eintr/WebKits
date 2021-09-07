package main

import (
	"flag"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/impact-eintr/WebKits/cmd/test/utils"
)

func main() {
	var ip string
	var port string
	var rounum int
	var ports []string
	var scanports []string
	var wgp sync.WaitGroup

	defaultports := [...]string{
		"21", "22", "23", "25", "80", "443", "8080",
		"110", "135", "139", "445", "389", "489", "587", "1433", "1434",
		"1521", "1522", "1723", "2121", "3306", "3389", "4899", "5631",
		"5632", "5800", "5900", "7071", "43958", "65500", "4444", "8888",
		"6789", "4848", "5985", "5986", "8081", "8089", "8443", "10000",
		"6379", "7001", "7002", "2049", "27017", "27018",
	}

	flag.StringVar(&ip, "u", "127.0.0.1", "扫描IP地址")
	flag.StringVar(&port, "p", "", "扫描的端口")
	flag.IntVar(&rounum, "n", 4, "协程数")
	flag.Parse()
	ips := utils.Headcheck(ip)
	if len(port) != 0 {
		ports = strings.Split(port, "-")
		startport, _ := strconv.Atoi(ports[0])
		endport, _ := strconv.Atoi(ports[1])
		for num := startport; num <= endport; num++ {
			scanports = append(scanports, strconv.Itoa(num))
		}
	} else {
		scanports = defaultports[:]
	}

	for _, v := range ips {
		scanchan := make(chan string, len(scanports))
		openchan := make(chan string, len(scanports))
		// 初始化数据源
		for _, value := range scanports {
			scanchan <- value
		}
		// 数据分发任务结束
		close(scanchan)

		start := time.Now()
		fmt.Println("任务开启")
		for i := 0; i < rounum; i++ {
			wgp.Add(1)
			go utils.CommandTcp(v, scanchan, openchan, &wgp)
		}

		wgp.Wait()

		// 任务已经结束 关掉数据接收
		close(openchan)
		end := time.Since(start)
		for {
			openport, ok := <-openchan
			if !ok {
				fmt.Println("扫描结束，无开放端口")
				break
			}
			fmt.Println("-------------------开放的端口----------------", ip, ":", openport)
		}
		fmt.Println("花费的时间", end)
	}

	fmt.Println("扫描结束")
}
