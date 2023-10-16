package main

import (
	"bufio"
	"errors"
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
	stdin := flag.Bool("i", false, "Use STDIN as input")
	file := flag.String("f", "", "Use file as input")
	flag.Parse()

	if *file == "" && !*stdin {
		flag.Usage()
		os.Exit(1)
	}

	input, err := selectInput(stdin, file)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer input.Close()

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

func selectInput(stdin *bool, file *string) (*os.File, error) {
	if *stdin {
		return os.Stdin, nil
	}

	if *file != "" {
		return os.Open(*file)
	}

	return nil, errors.New("unknown flag")
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
