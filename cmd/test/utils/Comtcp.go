package utils

import (
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

func CommandTcp(ip string, scanchan chan string, openchan chan string, wgscan *sync.WaitGroup) {
	defer func() {
		wgscan.Done()
	}()
	for {
		port, isend := <-scanchan
		if !isend {
			log.Println("没有数据了")
			break
		}
		fmt.Println("扫描端口:" + port)
		_, err := net.DialTimeout("tcp", ip+":"+port, time.Second)
		if err == nil {
			openchan <- port
		}
	}
}
