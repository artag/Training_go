package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	actions "rggo/fileSystem/ugz/actions"
)

func main() {
	src := flag.String("src", "", "Source root directory with archived files")
	dst := flag.String("dst", "", "Destination directory with extracted files")
	flag.Parse()

	if !actions.CheckFlagsAreNotEmpty(*src, *dst) {
		flag.Usage()
		os.Exit(1)
	}

	if err := run(*src, *dst, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func run(src, dst string, out io.Writer, outErr io.Writer) error {
	return filepath.Walk(
		src,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			ext := strings.Trim(filepath.Ext(path), ".")
			toFilter := actions.FilterInfo{
				IsDir:     info.IsDir(),
				Extension: ext,
			}
			exclude, err := actions.ExcludePath(toFilter)
			if err != nil {
				return err
			}

			if exclude {
				return nil
			}

			return actions.ExtractFile(src, dst, path, out)
		})
}
