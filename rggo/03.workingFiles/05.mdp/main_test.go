package main

import (
	"bytes"
	"io/ioutil" // To read data from files
	"os"
	"strings"
	"testing"
)

const (
	inputFile              = "./testdata/test1.md"
	goldenFile             = "./testdata/test1.md.html"
	goldenFileWithTemplate = "./testdata/test2.md.html"
	templateFile           = "./testdata/template.html.tmpl"
)

func TestParseContent(t *testing.T) {
	// Arrange
	env := resetEnvironment()
	defer restoreEnvironment(env)

	expected, err := ioutil.ReadFile(goldenFile)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	actual, err := parseContentFromFile(inputFile, "")
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assertBytes(expected, actual, t)
}

func TestRun(t *testing.T) {
	// Arrange
	env := resetEnvironment()
	defer restoreEnvironment(env)

	var mockStdOut bytes.Buffer

	// Act
	if err := run_file(inputFile, "", &mockStdOut, true); err != nil {
		t.Fatal(err)
	}

	resultFile := strings.TrimSpace(mockStdOut.String())
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

func TestParseContentWithTemplate(t *testing.T) {
	// Arrange
	env := resetEnvironment()
	defer restoreEnvironment(env)

	expected, err := ioutil.ReadFile(goldenFileWithTemplate)
	if err != nil {
		t.Fatal(err)
	}

	// Act
	actual, err := parseContentFromFile(inputFile, templateFile)
	if err != nil {
		t.Fatal(err)
	}

	// Assert
	assertBytes(expected, actual, t)
}

func TestRunWithTemplate(t *testing.T) {
	// Arrange
	env := resetEnvironment()
	defer restoreEnvironment(env)

	var mockStdOut bytes.Buffer

	// Act
	if err := run_file(inputFile, templateFile, &mockStdOut, true); err != nil {
		t.Fatal(err)
	}

	resultFile := strings.TrimSpace(mockStdOut.String())
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

	expected, err := ioutil.ReadFile(goldenFileWithTemplate)
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

func resetEnvironment() string {
	env := os.Getenv(EnvironmentVariable)
	if env == "" {
		return env
	}

	os.Setenv(EnvironmentVariable, "")
	return env
}

func restoreEnvironment(env string) {
	if env == "" {
		return
	}

	os.Setenv(EnvironmentVariable, env)
}
