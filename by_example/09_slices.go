package main

import "fmt"

func main() {
	s := make([]string, 3)
	fmt.Println("emp:", s)

	s[0] = "a"
	s[1] = "b"
	s[2] = "c"
	fmt.Println("set:", s)
	fmt.Println("get:", s[2])
	fmt.Println("len:", len(s))

	fmt.Println("------")

	s = append(s, "d")
	s = append(s, "e", "f")
	fmt.Println("apd:", s)

	fmt.Println("------")

	c := make([]string, len(s))
	copy(c, s)
	fmt.Println("cpy:", c)

	fmt.Println("------")

	c[0] = "a2"
	fmt.Println("arr1:", s)
	fmt.Println("arr2:", c)

	fmt.Println("------")

	l := s[2:5]
	fmt.Println("sl1:", l)

	l = s[:5]
	fmt.Println("sl2:", l)

	l = s[2:]
	fmt.Println("sl3:", l)

	fmt.Println("------")

	t := []string{"g", "h", "i"}
	fmt.Println("dcl:", t)

	fmt.Println("------")

	twoD := make([][]int, 3)
	for i := 0; i < 3; i++ {
		innerLen := i + 1
		twoD[i] = make([]int, innerLen)
		for j := 0; j < innerLen; j++ {
			twoD[i][j] = i + j
		}
	}
	fmt.Println("2d:", twoD)

	fmt.Println("------")

	d := []int{1, 2, 3, 4, 5}
	fmt.Println("dcl:", d)
	d2 := remove(d, 2)
	fmt.Println("del:", d2)
}

func remove(s []int, i int) []int {
	s1 := s[:i]
	s2 := s[i+1:]
	return append(s1, s2...)
}
