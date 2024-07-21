package base

import "github.com/billziss-gh/cgofuse/fuse"
import "fmt"

type BaseFS struct {
	fuse.FileSystemBase
}

func (fs *BaseFS) Mount(mountpoint string) error {
	host := fuse.NewFileSystemHost(fs)
	if !host.Mount(mountpoint, nil) {
		return fmt.Errorf("failed to mount filesystem at %s", mountpoint)
	}
	return nil
}
