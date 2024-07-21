package main

import (
	"fmt"
	"os"

	"github.com/1611Dhruv/file-systems/pkg/filesystem/vsf"
	"github.com/1611Dhruv/file-systems/pkg/tui"
)

func main() {
	if len(os.Args) > 1 {
		vsf.Run()
	} else {
		if err := tui.Start(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to start TUI: %v\n", err)
			os.Exit(1)
		}
	}
}
