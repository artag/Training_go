package main

import (
	"flag" // To handle command-line flags
	"fmt"  // To print formatted output
	"io"   // To use io.Writer interface
	"log"
	"os"            // To communicate wth the operating system
	"path/filepath" // To handle file paths
)

const (
	// Readable and writable by the owner but only readable by anyone else.
	FilePermissions = 0644
	DirPermissions  = 0755
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
	// log destination writer
	log io.Writer
	// archive directory
	archive string
	// from last modified datetime
	from string
}

func main() {
	// Parsing command line flags
	root := flag.String("root", ".", "Root directory to start")
	logFile := flag.String("log", "", "Log deletes to this file")
	// Action options
	archive := flag.String("archive", "", "Archive directory")
	list := flag.Bool("list", false, "List files only")
	del := flag.Bool("del", false, "Delete files")
	// Filter options
	ext := flag.String("ext", "", "File extension to filter out")
	size := flag.Int64("size", 0, "Minimum file size")
	from := flag.String("from", "", "Last modified datetime to filter out")
	flag.Parse()

	var (
		f   = os.Stdout
		err error
	)

	if *logFile != "" {
		// O_APPEND - Enables data to be appended to the end of file.
		// O_CREATE - Creates the file in case it doesn't exists.
		// O_RDWR - Opens the file for reading and writing.
		f, err = os.OpenFile(*logFile, os.O_APPEND|os.O_CREATE|os.O_RDWR, FilePermissions)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer f.Close()
	}

	c := config{
		ext:     *ext,
		size:    *size,
		list:    *list,
		del:     *del,
		log:     f,
		archive: *archive,
		from:    *from,
	}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(root string, out io.Writer, cfg config) error {
	logger := log.New(cfg.log, "DELETED FILE: ", log.LstdFlags)

	return filepath.Walk(root,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			filter, err := filterOut(path, cfg, info)
			if err != nil {
				return err
			}

			if filter {
				return nil
			}

			// List file
			if cfg.list {
				return listFile(path, out)
			}

			// Archive file
			if cfg.archive != "" {
				if err := archiveFile(cfg.archive, root, path); err != nil {
					return err
				}
			}

			// Delete file
			if cfg.del {
				return deleteFile(path, logger)
			}

			// List is the action by default
			return listFile(path, out)
		})
}
