package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// 创建子进程
	cmd := exec.Command("/bin/bash", "-c", `sleep 100 && echo "child: exit"`)

	cmd.Stdout = os.Stdout
	cmd.Start()

	fmt.Printf("parent id: %d ,child id %d \n", os.Getpid(), cmd.Process.Pid)
	//fmt.Println("parent: wake up")

	// 为子进程收尸
	//err := cmd.Wait()
	//if err != nil {
	//	panic(err)
	//}
}
