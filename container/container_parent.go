package container

import (
	"strings"
	"syscall"
	"os/exec"
	"os"
	"fmt"
)


// command =     /bin/sh
func NewParentProcess(tty bool, volume string) (*exec.Cmd, *os.File) {
	
	fmt.Println("NewParentProcess")
	readPipe, writePipe, err := NewPipe()
	if err != nil {
		fmt.Errorf("file pipe err: %v", err)
	}
	// readPipe, writePipe, err := NewPipe()
	// if err != nil {
	// 	return nil, nil
	// }
	initCmd, err := os.Readlink("/proc/self/exe")
	if err != nil {
		return nil, nil
	}

	// args := []string{"init", command}

	args := []string{"init"}
	
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
	// 从文件的readPipe 准备读入 command
	cmd.ExtraFiles = []*os.File{readPipe}
	mntURL := "/root/mnt"
	rootURL := "/root"
	NewWorkSpace(rootURL, mntURL, volume)
	// cmd.Dir = "/mnt/c/docker"
	cmd.Dir = mntURL
	return cmd, writePipe
}

func NewPipe() (*os.File, *os.File, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, nil, err
	}
	return read, write, nil
}

func NewWorkSpace(rootURL string, mntURL string, volumn string) {
	CreateReadOnlyLayer(rootURL)
	CreateWriteLayer(rootURL)
	CreateMountPoint(rootURL, mntURL)
	if (volumn != "") {
		fmt.Println("Volume not support !!")
	}
} 
func CreateReadOnlyLayer(rootURL string) {
	busyboxURL := rootURL + "/busybox"
	// busyboxTarURL := rootURL + "/busybox.tar"
	busyboxTarURL := "/mnt/c/docker/busybox.tar"

	exist, err := PathExists(busyboxURL)
	if err != nil {
		fmt.Errorf("error : %v", err)
	}
	if exist == false {
		if err := os.Mkdir(busyboxURL, 0777); err != nil {
			fmt.Println("Mkdir busybox dir error")
		}
		if _, err := exec.Command("tar", "-xvf", busyboxTarURL, "-C", busyboxURL).CombinedOutput(); err != nil {
			fmt.Println("Untar dir error")
		}
	}

}

func CreateWriteLayer(rootURL string) {
	writeURL := rootURL + "/writeLayer"
	if err := os.Mkdir(writeURL, 0777); err != nil {
		fmt.Println("Mkdir writeLayer error")
	}
}

func CreateMountPoint(rootURL string, mntURL string) {
	if err := os.Mkdir(mntURL, 0777); err != nil {
		fmt.Println("Mkdir mntURL error")
	}
	// 参考这个参数
	// http://man7.org/linux/man-pages/man2/mount.2.html
	dirs := "dirs=" + rootURL + "/writeLayer:" + rootURL + "/busybox"
	cmd := exec.Command("mount", "-t", "aufs", "-o", dirs, "none", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Println(err)
	}
}
// delete the aufs filesystem
func DeleteWorkSpace(rootURL string, mntURL string, volume string) {
	if volume != "" {
		volumeURLs := volumeUrlExtract(volume)
		length := len(volumeURLs)
		if length == 2 && volumeURLs[0] != "" && volumeURLs[1] != "" {

		} else {
			DeleteMountPoint(rootURL, mntURL)
		}
	} else {
		DeleteMountPoint(rootURL, mntURL)
	}

}

func DeleteMountPoint(rootURL string, mntURL string) {
	cmd := exec.Command("umount", mntURL)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Errorf("%v", err)
	}
	if err := os.RemoveAll(mntURL); err != nil {
		fmt.Errorf("Remove mountpoint dir %s error %v", mntURL, err)
	}
}

func DeleteWriteLayer(rootURL string) {
	writeURL := rootURL + "/writeLayer"
	if err := os.RemoveAll(writeURL); err != nil {
		fmt.Errorf("Remove writeLayer dir %s error %v", writeURL, err)
	}
}

func volumeUrlExtract(volume string) ([]string) {
	var volumeURLs []string
	volumeURLs = strings.Split(volume, ":")
	return volumeURLs
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
