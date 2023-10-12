package main

import (
	"bytes"
	"io/ioutil" // To read data from files
	"os"
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
	assertBytes(expected, actual, t)
}

func TestRun(t *testing.T) {
	// Act
	if err := run(inputFile); err != nil {
		t.Fatal(err)
	}

	exists, err := exists(resultFile)
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		defer os.Remove(resultFile)
	}

	actual, err := ioutil.ReadFile(resultFile)
	if err != nil {
		t.Fatal(err)
	}

	expected, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assertBytes(expected, actual, t)
}

func assertBytes(expected []byte, actual []byte, t *testing.T) {
	if bytes.Equal(expected, actual) {
		return
	}

	t.Logf("expected:\n%s\n", expected)
	t.Logf("result:\n%s\n", actual)
	t.Error("Actual content does not match expected")
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}
