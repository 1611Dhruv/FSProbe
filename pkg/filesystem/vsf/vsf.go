package vsf

import (
	"encoding/binary"
	"log"
	"os"

	"github.com/1611Dhruv/file-systems/internal/config"
	"github.com/billziss-gh/cgofuse/fuse"
)

type VSF struct {
}

// VSF_FILE represents the file system
type VSF_FILE struct {
	fuse.FileSystemBase
	superBlock *VSF_SUPERBLOCK
	inodes     []VSF_INODE
	dirents    []VSF_DIRENT
}

type VSF_SUPERBLOCK struct {
	magic      [config.MagicLength]byte
	iNodeCount int
	dataCount  int
	imap       []byte
	dmap       []byte
}

type VSF_INODE struct {
	ino    uint64
	stat   fuse.Stat_t
	blocks []*byte
}

type VSF_DIRENT struct {
	name [256]byte
	ino  uint64
}

// MakeVSF creates a new VSF file system
func MakeVSF(iNodeCount, dataCount int, filePath string) {

	// Round up dataCount and iNodeCount to be aligned
	if dataCount%config.Alignment != 0 {
		dataCount += config.Alignment - dataCount%config.Alignment
	}
	if iNodeCount%config.Alignment != 0 {
		iNodeCount += config.Alignment - iNodeCount%config.Alignment
	}

	// Open the file for reading and create it if it doesn't exist
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	superBlock := new(VSF_SUPERBLOCK)
	for i := 0; i < config.MagicLength; i++ {
		superBlock.magic[i] = config.VSFMagic[i]
	}

	superBlock.imap = make([]byte, config.BlockSize*config.Alignment)
	superBlock.dmap = make([]byte, config.BlockSize*config.Alignment)
	superBlock.iNodeCount = iNodeCount
	superBlock.dataCount = dataCount

	err = binary.Write(file, binary.LittleEndian, superBlock)
	// Write out the superBlock
	if err != nil {
		log.Fatalf("failed to write to file: %v", err)
	}

	// Write out the remaining bytes
	bytes := make([]byte, config.BlockSize*(iNodeCount+dataCount))
	_, err = file.WriteAt(bytes, config.BlockSize)
	if err != nil {
		log.Fatalf("failed to write to file: %v", err)
	}
	file.Close()
}

func InitVSF(imageFile string) *VSF_FILE {
	file, err := os.OpenFile(imageFile, os.O_RDWR, 0755)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}

	superBlock := new(VSF_SUPERBLOCK)
	err = binary.Read(file, binary.LittleEndian, superBlock)
	if err != nil {
		log.Fatalf("failed to read from file: %v", err)
	}

	fs := new(VSF_FILE)
	fs.superBlock = superBlock

	// Make sure we have enough memory to hold all the inodes and dirents
	fs.inodes = make([]VSF_INODE, superBlock.iNodeCount)
	fs.dirents = make([]VSF_DIRENT, superBlock.dataCount)

	return fs
}

func (fs *VSF) Init() {
	// Initialize filesystem (e.g., read disk image, setup in-memory structures)
}

func (fs *VSF) Access(path string, mode uint32) (errc int) {
	log.Printf("Access(path=%s, mode=%d)", path, mode)
	return 0
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
	fs := InitVSF("disk.img")
	host := fuse.NewFileSystemHost(*fs)
	host.Mount("", os.Args[1:])
}
