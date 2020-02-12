package subsystem

import (
	"io/ioutil"
	"fmt"
	"path"
)

type MemorySubSystem struct {

}

func (s *MemorySubSystem) Set(cgroupPath string, resc *ResourceConfig) error {
	if subsysCgroupPath, err := GetCgroupPath(s.Name(), cgroupPath, true); err == nil {
		if resc.MemoryLimit != "" {
			if err := ioutil.WriteFile(path.Join(subsysCgroupPath, "memory.limit_in_bytes"), []byte(resc.MemorySubSystem)); err != nil {
				return fmt.Errorf("set cgroup memory fail %v", err)
			}
		}
		return nil
	} else {
		return err
	}
}