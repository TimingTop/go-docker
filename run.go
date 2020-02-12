package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/go-docker/container"
	"os"
)

func Run(tty bool, command string) {
	containerID := genUUID()
	fmt.Println("ContainerId: " + containerID)
	// 创建容器的父进程
	parent, _ := container.NewParentProcess(tty, command)
	fmt.Println(parent)
	if err := parent.Start(); err != nil {
		fmt.Println("error!")
		fmt.Println(err)
		
		os.Exit(-1)
	}
	parent.Wait()
	os.Exit(-1)
}

func genUUID() string {
	id := uuid.New()
	return id.String()
}

