package vsf

import (
	"log"
	"os"

	"github.com/billziss-gh/cgofuse/fuse"
)

type VSF struct {
	fuse.FileSystemBase
}

func NewVSF() *VSF {
	return &VSF{}
}

func (fs *VSF) Init() {
	// Initialize filesystem (e.g., read disk image, setup in-memory structures)
}

func (fs *VSF) Open(path string, flags int) (errc int, fh uint64) {
	log.Printf("Open(path=%s, flags=%d)", path, flags)
	return 0, 0
}

func (fs *VSF) Read(path string, buff []byte, ofst int64, fh uint64) (n int) {
	log.Printf("Read(path=%s, size=%d, ofst=%d)", path, len(buff), ofst)
	return 0
}

func (fs *VSF) Write(path string, buff []byte, ofst int64, fh uint64) (n int) {
	log.Printf("Write(path=%s, size=%d, ofst=%d)", path, len(buff), ofst)
	return len(buff)
}

func (fs *VSF) Getattr(path string, stat *fuse.Stat_t, fh uint64) (errc int) {
	log.Printf("Getattr(path=%s)", path)
	if path == "/" {
		stat.Mode = fuse.S_IFDIR | 0755
	} else {
		stat.Mode = fuse.S_IFREG | 0644
		stat.Size = 0
	}
	return 0
}

func Run() {
	fs := NewVSF()
	host := fuse.NewFileSystemHost(fs)
	host.Mount("", os.Args[1:])
}
