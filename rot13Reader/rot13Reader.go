package main

import (
	"io"
	"os"
	"strings"
)

type rot13Reader struct {
	r io.Reader
}

func (rot13 rot13Reader) Read(b []byte) (n int, err error) {
	length, strEnd := rot13.r.Read(b)
	for i := 0; i < length; i++ {
		b[i] = decode(b[i])
	}

	return length, strEnd
}

func decode(s byte) byte {
	if (64 < s && s < 78) ||
		(96 < s && s < 110) {
		return s + 13
	}
	if (77 < s && s < 91) ||
		(109 < s && s < 123) {
		return s - 13
	}

	return s
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}
