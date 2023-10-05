package main

import (
	"bufio" // read text
	"fmt"   // print formatted output
	"io"    // provides th io.Reader interface
	"os"    // use operating system resources
)

func main() {
	input := os.Stdin
	result := count(input)
	fmt.Println(result)
}

func count(r io.Reader) int {
	// A scanner is used to read text from a Reader (such as files)
	scanner := bufio.NewScanner(r)

	// Define the scanner split type to words (default is split by lines)
	scanner.Split(bufio.ScanWords)

	// Defining a counter
	wc := 0

	// For every word scanned, increment the counter
	for scanner.Scan() {
		wc++
	}

	// Return the total
	return wc
}
