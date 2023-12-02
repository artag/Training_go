package add

import "fmt"

func add(a, b int) int {
	callsome()
	return a + b
}

func callsome() error {
	fmt.Println("run add")
	return nil
}
