package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	z := 1.0
	eps := 0.001
	diff := 1.0
	for diff > eps {
		z -= (z * z - x) / (2 * z)
		diff = z * z - x
		z = z + 0.00001
	}
	
	return z
}

func main() {
	i := float64(5)
	fmt.Println("My sqrt impl:", Sqrt(i))
	fmt.Println("Standard lib:", math.Sqrt(i))
}
