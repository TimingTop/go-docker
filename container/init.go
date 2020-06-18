package container

import (
	"path/filepath"
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

	//defaultMountFlags := syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV
    // http://man7.org/linux/man-pages/man2/mount.2.html 
	//syscall.Mount("proc", "/proc", "proc", uintptr(defaultMountFlags), "")

	setUpMount()
	argv := []string{command}

	if err := syscall.Exec(command, argv, os.Environ()); err != nil {
		fmt.Println(err)
	}
	return nil
}

func setUpMount() {
	pwd, err := os.Getwd()
	fmt.Println("pwd = " + pwd)
	if err != nil {
		fmt.Println(err)
		return
	}
	pivotRoot(pwd)

	syscall.Mount("proc", "/proc", "proc", syscall.MS_NOEXEC | syscall.MS_NOSUID | syscall.MS_NODEV, "")
	syscall.Mount("tmpfs", "/dev", "tmpfs", syscall.MS_NOSUID|syscall.MS_STRICTATIME, "mode=755")

}

func pivotRoot(root string) error {
	err := syscall.Mount(root, root, "bind", syscall.MS_BIND|syscall.MS_REC, "")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	pivotDir := filepath.Join(root, ".pivot_root")
	err = os.Mkdir(pivotDir, 0777)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	err = syscall.PivotRoot(root, pivotDir)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	err = syscall.Chdir("/")
	if err != nil {
		fmt.Println(err)
		return nil
	}

	pivotDir = filepath.Join("/", ".pivot_root")
	err = syscall.Unmount(pivotDir, syscall.MNT_DETACH)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return os.Remove(pivotDir)
	
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
