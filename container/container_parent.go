package container

import (
	"syscall"
	"os/exec"
	"os"
	"fmt"
)


// command =     /bin/sh
func NewParentProcess(tty bool, command string) (*exec.Cmd, *os.File) {
	fmt.Println("NewParentProcess")
	// readPipe, writePipe, err := NewPipe()
	// if err != nil {
	// 	return nil, nil
	// }
	initCmd, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return nil, nil
	}

	args := []string{"init", command}
	
	cmd := exec.Command(initCmd, args...)
	// syscall.CLONE_NEWNET   not support 
	cmd.SysProcAttr = &syscall.SysProcAttr {
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS | syscall.CLONE_NEWIPC,
	}

	if tty {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	// cmd.ExtraFiles = []*os.File{readPipe}

	return cmd, nil


}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}