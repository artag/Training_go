package main

import (
	"os"
	"testing"
)

func TestFilterOut(t *testing.T) {
	testCases := []struct {
		name     string
		file     string
		ext      string
		minSize  int64
		expected bool
	}{
		{"FilerNoExtension", "testdata/dir.log", "", 0, false},
		{"FilerExtensionMatch", "testdata/dir.log", ".log", 0, false},
		{"FilerExtensionNoMatch", "testdata/dir.log", ".sh", 0, true},
		{"FilerExtensionSizeMatch", "testdata/dir.log", ".log", 10, false},
		{"FilerExtensionSizeNoMatch", "testdata/dir.log", ".log", 20, true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			info, err := os.Stat(tc.file)
			if err != nil {
				t.Fatal(err)
			}

			f := filterOut(tc.file, tc.ext, tc.minSize, info)
			if f != tc.expected {
				t.Errorf("Expected: '%t'\nActual '%t'", tc.expected, f)
			}
		})
	}
}
