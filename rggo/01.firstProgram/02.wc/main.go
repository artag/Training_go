package main

import (
	"bufio"
	"flag" // use to create and manage command line flags
	"fmt"
	"io"
	"os"
)

func main() {
	// Defining a boolean flag -l to count lines instead of words
	lines := flag.Bool("l", false, "Count lines")
	// Parsing the flags provided by the user
	flag.Parse()

	input := os.Stdin
	result := count(input, *lines)
	fmt.Println(result)
}

func count(r io.Reader, countLines bool) int {
	scanner := bufio.NewScanner(r)

	// If the count lines flag is not set, we want to count words so we define
	// the scanner split type to words (default is split by lines)
	if !countLines {
		scanner.Split(bufio.ScanWords)
	}

	wc := 0
	for scanner.Scan() {
		wc++
	}

	return wc
}
