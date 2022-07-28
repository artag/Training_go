package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	pln := fmt.Println

	p := filepath.Join("dir1", "dir2", "filename")
	pln("p =", p)

	pln(filepath.Join("dir1//", "filename"))
	pln(filepath.Join("dir1/../dir1", "filename"))

	pln("Dir(p):", filepath.Dir(p))
	pln("Base(p):", filepath.Base(p))
	dir, base := filepath.Split(p)
	pln("After Split(p). dir:", dir, ", base:", base)

	pln("Is absolute 'dir/file' path:", filepath.IsAbs("dir/file"))
	pln("Is absolute '/dir/file' path:", filepath.IsAbs("/dir/file"))

	filename := "config.json"
	ext := filepath.Ext(filename)
	pln("ext:", ext)
	pln("Trim suffix 'json' from 'config.json':", strings.TrimSuffix(filename, ext))

	pln("---")

	rel, err := filepath.Rel("a/b", "a/b/t/file")
	check(err)
	pln("Rel:", rel)

	rel, err = filepath.Rel("a/b", "a/c/t/file")
	check(err)
	pln("Rel:", rel)
}
