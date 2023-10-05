package main

import (
	"bufio"
	"flag" // use to create and manage command line flags
	"fmt"
	"io"
	"os"
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
	flag.Parse()

	input := os.Stdin

	var result = 0
	if *bytes {
		result = count(input, countBytes)
	} else if *lines {
		result = count(input, countLines)
	} else {
		result = count(input, countWords)
	}

	fmt.Println(result)
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
