package main

import (
	"flag"          // To handle command-line flags
	"fmt"           // To print formatted output
	"io"            // To use io.Writer interface
	"os"            // To communicate wth the operating system
	"path/filepath" // To handle file paths
)

type config struct {
	// extension to filter out
	ext string
	// minimal file size
	size int64
	// list files
	list bool
	// delete files
	del bool
}

func main() {
	// Parsing command line flags
	root := flag.String("root", ".", "Root directory to start")
	// Action options
	list := flag.Bool("list", false, "List files only")
	del := flag.Bool("del", false, "Delete files")
	// Filter options
	ext := flag.String("ext", "", "File extension to filter out")
	size := flag.Int64("size", 0, "Minimum file size")
	flag.Parse()

	c := config{
		ext:  *ext,
		size: *size,
		list: *list,
		del:  *del,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root string, out io.Writer, cfg config) error {
	return filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if filterOut(path, cfg.ext, cfg.size, info) {
				return nil
			}

			// List file
			if cfg.list {
				return listFile(path, out)
			}

			// Delete file
			if cfg.del {
				return deleteFile(path)
			}

			// List is the action by default
			return listFile(path, out)
		})
}
