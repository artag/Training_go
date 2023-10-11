package main

import (
	"bytes"
	"io/ioutil" // To read data from files
	"testing"
)

const (
	inputFile  = "./testdata/test1.md"
	goldenFile = "./testdata/test1.md.html"
	resultFile = "test1.md.html"
)

func TestParseContent(t *testing.T) {
	// Arrange
	input, err := ioutil.ReadFile(inputFile)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	actual := parseContent(input)

	// Assert
	if !bytes.Equal(expected, actual) {
		t.Logf("expected:\n%s\n", expected)
		t.Logf("result:\n%s\n", actual)
		t.Error("Actual content does not match expected")
	}
}
