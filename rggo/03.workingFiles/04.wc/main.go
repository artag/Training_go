package main

import (
	"bufio"
	"errors"
	"flag" // use to create and manage command line flags
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	countWords = iota
	countLines = iota
	countBytes = iota
)

func main() {
	// Defining a boolean flag -b to count bytes
	bytes := flag.Bool("b", false, "Count bytes")
	lines := flag.Bool("l", false, "Count lines")
	stdin := flag.Bool("i", false, "Use STDIN as input")
	file := flag.String("f", "", "Use file as input")
	flag.Parse()

	if *file == "" && !*stdin {
		flag.Usage()
		os.Exit(1)
	}

	// Select input and count
	var result int = 0
	if *stdin {
		result = countDataStream(os.Stdin, bytes, lines)
	} else if *file != "" {
		split := strings.Split(*file, " ")
		result = countMultipleFiles(split, bytes, lines)
	} else {
		err := errors.New("unknown work mode")
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(result)
}

func countMultipleFiles(files []string, bytes *bool, lines *bool) int {
	var result int = 0
	for _, f := range files {
		input, err := os.Open(f)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			os.Exit(1)
		}
		defer input.Close()
		result += countDataStream(input, bytes, lines)
	}
	return result
}

func countDataStream(input *os.File, bytes *bool, lines *bool) int {
	var result = 0
	if *bytes {
		result = count(input, countBytes)
	} else if *lines {
		result = count(input, countLines)
	} else {
		result = count(input, countWords)
	}

	return result
}

func count(r io.Reader, countOptions int) int {
	scanner := bufio.NewScanner(r)

	if countOptions == countBytes {
		scanner.Split(bufio.ScanBytes)
	}

	if countOptions == countWords {
		scanner.Split(bufio.ScanWords)
	}

	wc := 0
	for scanner.Scan() {
		wc++
	}

	return wc
}
