package container

import (
	"fmt"
	"syscall"
	// "os/exec"
	"strings"
	"io/ioutil"
	"os"
)


// command =  /bin/sh
func RunContainerInitProcess(command string, args []string) error {
	// cmdArray := readUserCommand()
	fmt.Println("RunContainerInitProcess: " + command)

	// if cmdArray == nil || len(cmdArray) == 0 {
	// 	return fmt.Errorf("CmdArray is nil")
	// }

	// path, err := exec.LookPath(cmdArray[0])
	// if err != nil {
	// 	return err
	// }
	// if err := syscall.Exec(path, cmdArray[0:], os.Environ()); err != nil {
	// 	fmt.Printf("syscall is null")
	// }
	// return nil

	defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
    // http://man7.org/linux/man-pages/man2/mount.2.html 
	syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	argv := []string{command}

	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		fmt.Println(err)
	}
	return nil
}

func readUserCommand() []string {
	pipe := os.NewFile(uintptr(3), "pipe")
	defer pipe.Close()
	msg, err := ioutil.ReadAll(pipe)
	if err != nil {
		
		return nil
	}
	msgStr := string(msg)
	return strings.Split(msgStr, " ")
}
