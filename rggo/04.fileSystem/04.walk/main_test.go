package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func TestRunArchive(t *testing.T) {
	testCases := []struct {
		name         string
		cfg          config
		extNoArchive string
		nArchive     int
		nNoArchive   int
	}{
		{
			name:         "ArchiveExtensionNoMatch",
			cfg:          config{ext: ".log"},
			extNoArchive: ".gz",
			nArchive:     0,
			nNoArchive:   10},
		{
			name:         "ArchiveExtensionMatch",
			cfg:          config{ext: ".log"},
			extNoArchive: "",
			nArchive:     10,
			nNoArchive:   0},
		{
			name:         "ArchiveExtensionMixed",
			cfg:          config{ext: ".log"},
			extNoArchive: ".gz",
			nArchive:     5,
			nNoArchive:   5},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Buffer for RunArchuve output
			var buffer bytes.Buffer

			// Create temp dirs for RunArchive test
			tempDir, cleanup := createTempDir(t, map[string]int{
				tc.cfg.ext:      tc.nArchive,
				tc.extNoArchive: tc.nNoArchive,
			})
			defer cleanup()

			archiveDir, cleanupArchive := createTempDir(t, nil) // nil - we don't need any files
			defer cleanupArchive()

			// Act
			tc.cfg.archive = archiveDir
			if err := run(tempDir, &buffer, tc.cfg); err != nil {
				t.Fatal(err)
			}

			// Assert
			pattern := filepath.Join(tempDir, fmt.Sprintf("*%s", tc.cfg.ext))
			expFiles, err := filepath.Glob(pattern) // Glob - find all file names from directory
			if err != nil {
				t.Fatal(err)
			}

			expOut := strings.Join(expFiles, "\n")
			actOut := strings.TrimSpace(buffer.String())
			if expOut != actOut {
				t.Errorf("Expected: %q\nActual: %q\n", expOut, actOut)
			}

			// Validate the number of files archived
			filesArchived, err := ioutil.ReadDir(archiveDir)
			if err != nil {
				t.Fatal(err)
			}

			if len(filesArchived) != tc.nArchive {
				t.Errorf("Exp ected: %d archived files\nActual: %d archived files\n",
					tc.nArchive, len(filesArchived))
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
			var (
				buffer    bytes.Buffer
				logBuffer bytes.Buffer
			)
			tc.cfg.log = &logBuffer
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

			// Assert logged lines
			expLogLines := tc.nDelete + 1
			lines := bytes.Split(logBuffer.Bytes(), []byte("\n"))
			if len(lines) != expLogLines {
				t.Errorf("Expected: %d log lines\nActual: %d log lines\n", expLogLines, len(lines))
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
			if err := ioutil.WriteFile(fpath, []byte("dummy"), FilePermissions); err != nil {
				t.Fatal(err)
			}
			//fmt.Printf("Create temp file: %s\n", fpath) // To debug
		}
	}

	return tempDir, func() { os.RemoveAll(tempDir) }
}
