package main

import "fmt"

func main() {
	m := make(map[string]int)

	m["k1"] = 7
	m["k2"] = 13
	fmt.Println("map:", m)
	m["k2"] = 14
	fmt.Println("map:", m)

	fmt.Println("---")

	v1 := m["k1"]
	fmt.Println("v1:", v1)
	v1 = m["k3"]
	fmt.Println("v1:", v1)

	fmt.Println("------")

	fmt.Println("len:", len(m))

	fmt.Println("------")

	delete(m, "k2")
	fmt.Println("map:", m)

	_, prs := m["k2"]
	fmt.Println("prs:", prs)

	fmt.Println("------")

	n := map[string]int{"foo": 1, "bar": 2}
	fmt.Println("map:", n)
}
