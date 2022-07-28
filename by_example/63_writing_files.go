package main

import (
	"bufio"
	"fmt"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	p := fmt.Println
	pf := fmt.Printf

	filename1 := "dat1"
	pf("Write to file \"%s\" ...\n", filename1)

	d1 := []byte("hello\ngo\n")
	err := os.WriteFile(filename1, d1, 0644)
	check(err)

	p("Done")
	printFileContent(filename1)

	//-----------------------------------

	filename2 := "dat2"
	pf("Open file \"%s\" to write ...\n", filename2)

	f, err := os.Create(filename2)
	check(err)
	defer f.Close()

	d2 := []byte{115, 111, 109, 101, 10}
	n2, err := f.Write(d2)
	check(err)
	fmt.Printf("wrote %d bytes\n", n2)

	n3, err := f.WriteString("writes\n")
	check(err)
	fmt.Printf("wrote %d bytes\n", n3)
	f.Sync()

	w := bufio.NewWriter(f)
	n4, err := w.WriteString("buffered\n")
	check(err)
	fmt.Printf("wrote %d bytes\n", n4)
	w.Flush()

	f.Close()
	p("Done")
	printFileContent(filename2)
}

func printFileContent(fname string) {
	dat, err := os.ReadFile(fname)
	check(err)
	fmt.Printf("File \"%s\":\n", fname)
	fmt.Print(string(dat))
	fmt.Println("---")
	fmt.Println()
}
