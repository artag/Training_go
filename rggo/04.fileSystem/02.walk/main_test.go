package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

const (
	// Readable and writable by the owner but only readable by anyone else.
	filePermissions = 0644
)

func TestRun(t *testing.T) {
	testCases := []struct {
		name     string
		root     string
		cfg      config
		expected string
	}{
		{
			name:     "NoFilter",
			root:     "testdata",
			cfg:      config{ext: "", size: 0, list: true},
			expected: "testdata/dir.log\ntestdata/dir2/script.txt\n"},
		{
			name:     "FilterExtensionMatch",
			root:     "testdata",
			cfg:      config{ext: ".log", size: 0, list: true},
			expected: "testdata/dir.log\n"},
		{
			name:     "FilterExtensionSizeMatch",
			root:     "testdata",
			cfg:      config{ext: ".log", size: 10, list: true},
			expected: "testdata/dir.log\n"},
		{
			name:     "FilterExtensionSizeNoMatch",
			root:     "testdata",
			cfg:      config{ext: ".log", size: 20, list: true},
			expected: ""},
		{
			name:     "FilterExtensionNoMatch",
			root:     "testdata",
			cfg:      config{ext: ".gz", size: 0, list: true},
			expected: ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			if err := run(tc.root, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			actual := buffer.String()
			if tc.expected != actual {
				t.Errorf("Expected: %q\nActual: %q\n", tc.expected, actual)
			}
		})
	}
}

func TestRunDelExtension(t *testing.T) {
	testCases := []struct {
		// Наименование теста
		name string
		cfg  config
		// Расширение файлов, которые не будут удалены
		extNoDelete string
		// Количество удаляемых файлов
		nDelete int
		// Количество оставшихся файлов
		nNoDelete int
		// Ожидаемый вывод на консоль
		expected string
	}{
		{
			// Create 10 files: file1.gz ... file10.gz
			// Останется 10 файлов, 0 будет удалено
			name:        "DeleteExtensionNoMatch",
			cfg:         config{ext: ".log", del: true},
			extNoDelete: ".gz",
			nDelete:     0,
			nNoDelete:   10,
			expected:    ""},
		{
			// Create 10 files: file1.log ... file10.log
			// Останется 0 файлов, 10 будет удалено
			name:        "DeleteExtensionMatch",
			cfg:         config{ext: ".log", del: true},
			extNoDelete: "",
			nDelete:     10,
			nNoDelete:   0,
			expected:    ""},
		{
			// Create 10 files: file1.log ... file5.log, file1.gz ... file5.gz
			// Останется 5 файлов, 5 будет удалено
			name:        "DeleteExtensionMixed",
			cfg:         config{ext: ".log", del: true},
			extNoDelete: ".gz",
			nDelete:     5,
			nNoDelete:   5,
			expected:    ""},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var buffer bytes.Buffer
			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.cfg.ext:     tc.nDelete,
				tc.extNoDelete: tc.nNoDelete,
			})
			defer cleanup()

			if err := run(tempDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			actual := buffer.String()
			if tc.expected != actual {
				t.Errorf("Expected: %q\nActual: %q\n", tc.expected, actual)
			}

			filesLeft, err := ioutil.ReadDir(tempDir)
			if err != nil {
				t.Error(err)
			}

			//fmt.Printf("Test: %s. Expected files left: %d, actual: %d\n", tc.name, len(filesLeft), tc.nNoDelete)
			if len(filesLeft) != tc.nNoDelete {
				t.Errorf("Expected: %d files\nActual: %d files\n", tc.nNoDelete, len(filesLeft))
			}
		})
	}
}

func createTempDir(t *testing.T, files map[string]int) (dirname string, cleanup func()) {
	t.Helper()

	tempDir, err := ioutil.TempDir("", "walktest")
	if err != nil {
		t.Fatal(err)
	}

	for k, n := range files {
		for j := 1; j <= n; j++ {
			fname := fmt.Sprintf("file%d%s", j, k)
			fpath := filepath.Join(tempDir, fname)
			if err := ioutil.WriteFile(fpath, []byte("dummy"), filePermissions); err != nil {
				t.Fatal(err)
			}
			// fmt.Printf("Create temp file: %s\n", fpath)   // To debug
		}
	}

	return tempDir, func() { os.RemoveAll(tempDir) }
}
