package main

import "fmt"

type rect struct {
	width, height int
}

func (r *rect) area() int {
	return r.width * r.height
}

func (r rect) perim() int {
	return 2 * (r.width + r.height)
}

func main() {
	r := rect{width: 10, height: 5}

	print("r", &r)

	fmt.Println("---------")

	rcp := r
	rcp.width = 6
	print("r", &r)
	print("rcp", &rcp)

	fmt.Println("---------")

	rcp2 := &r
	rcp2.width = 7
	print("r", &r)
	print("rcp", &rcp)
	print("rcp2", rcp2)
}

func print(name string, r *rect) {
	fmt.Println(name, "rect:", *r)
	fmt.Println(name, "area:", r.area())
	fmt.Println(name, "perim:", r.perim())
}
