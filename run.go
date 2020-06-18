package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-docker/container"
	"github.com/google/uuid"
)
// comArray = ["imageName", "/bin/ls", "-al"]
func Run(tty bool, comArray []string, volume string) {
	containerID := genUUID()
	fmt.Println("ContainerId: " + containerID)
	// 创建容器的父进程
	// 从 writePipe 写入 command    /bin/ls -al
	parent, writePipe:= container.NewParentProcess(tty, volume)
	fmt.Println(parent)
	if err := parent.Start(); err != nil {
		fmt.Println("error!")
		fmt.Println(err)
		
		os.Exit(-1)
	}
	// 写入命令
	sendInitCommand(comArray, writePipe)
	parent.Wait()

	mntURL := "/root/mnt"
	rootURL := "/root"
	container.DeleteWorkSpace(rootURL, mntURL, volume)
	os.Exit(0)
}

func sendInitCommand(comArray []string, writePipe *os.File) {
	command := strings.Join(comArray, " ")
	fmt.Printf("command all is %s", command)
	writePipe.WriteString(command)
	writePipe.Close()
}

func genUUID() string {
	id := uuid.New()
	return id.String()
}

