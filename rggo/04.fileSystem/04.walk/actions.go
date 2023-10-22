package main

import (
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func filterOut(path, ext string, minSize int64, info os.FileInfo) bool {
	if info.IsDir() || info.Size() < minSize {
		return true
	}

	if ext != "" && filepath.Ext(path) != ext {
		return true
	}

	return false
}

func archiveFile(dstDir, root, path string) error {
	info, err := os.Stat(dstDir)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", dstDir)
	}

	// Get archive filename wih path
	relDir, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil {
		return err
	}

	archiveFilename := filepath.Base(path)
	dst := fmt.Sprintf("%s.gz", archiveFilename)
	targetPath := filepath.Join(dstDir, relDir, dst)

	// Make directory. Or do nothing if directory already exists
	if err := os.MkdirAll(filepath.Dir(targetPath), DirPermissions); err != nil {
		return err
	}

	// Create and open destination file
	out, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, FilePermissions)
	if err != nil {
		return err
	}
	defer out.Close()

	// Open source file
	in, err := os.Open(path)
	if err != nil {
		return err
	}
	defer in.Close()

	// Compress input file
	zw := gzip.NewWriter(out)
	zw.Name = archiveFilename
	if _, err = io.Copy(zw, in); err != nil {
		return err
	}

	if err := zw.Close(); err != nil {
		return err
	}

	return out.Close()
}

func listFile(path string, out io.Writer) error {
	_, err := fmt.Fprintln(out, path)
	return err
}

func deleteFile(path string, log *log.Logger) error {
	if err := os.Remove(path); err != nil {
		return err
	}

	log.Println(path)
	return nil
}
