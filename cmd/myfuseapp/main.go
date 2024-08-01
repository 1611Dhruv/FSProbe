package main

import (
	"github.com/1611Dhruv/file-systems/pkg/filesystem/vsf"
)

func main() {
	vsf.InitializeVSF(100, 100, "disk.img")
}
