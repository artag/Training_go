package actions

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const (
	supportedExtension = "gz"
	filePermissions    = 0644
	dirPermissions     = 0755
)

func CheckFlagsAreNotEmpty(src, dst string) bool {
	if src != "" && dst != "" {
		return true
	}

	return false
}

type FilterInfo struct {
	IsDir     bool
	Extension string
}

func ExcludePath(info FilterInfo) (bool, error) {
	if info.IsDir {
		return true, nil
	}

	ext := strings.ToLower(info.Extension)
	if ext != supportedExtension {
		return true, nil
	}

	return false, nil
}

func ExtractFile(src, dst, path string, printer io.Writer) error {
	// Check destination
	info, err := os.Stat(dst)
	if err != nil {
		return err
	}

	if !info.IsDir() {
		return fmt.Errorf("%s is not a directory", dst)
	}

	// Get relative directory
	relDir, err := filepath.Rel(src, filepath.Dir(path))
	if err != nil {
		return err
	}

	// Get destination file path
	archiveFilename := filepath.Base(path)
	extractedFilename := GetFileWithoutGZ(archiveFilename)
	targetPath := filepath.Join(dst, relDir, extractedFilename)

	// Make target directory. Or do nothing if directory already exists
	targetDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(targetDir, dirPermissions); err != nil {
		return err
	}

	// Create and open destination file
	out, err := os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, filePermissions)
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

	// Extract input file
	unzipped, err := gzip.NewReader(in)
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, unzipped); err != nil {
		return err
	}

	// Print message
	msg := fmt.Sprintf("Extracted file '%s' to '%s'\n", path, targetPath)
	printer.Write([]byte(msg))

	return out.Close()
}

func GetFileWithoutGZ(filename string) string {
	suffixes := []string{".gz", ".GZ", ".Gz", ".gZ"}
	for _, sfx := range suffixes {
		if strings.HasSuffix(filename, sfx) {
			return strings.TrimSuffix(filename, sfx)
		}
	}

	return filename
}
