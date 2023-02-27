package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

var count int = 0

func makeRequest(ctx context.Context) (string, error) {
	count++

	if count <= 3 {
		return "", errors.New("some error happens")
	} else {
		return "some value", nil
	}
}

func main() {
	testBreaker()
}

func testBreaker() {
	fmt.Println("========== Breaker")
	fmt.Println()

	ctx := context.Background()
	breaker := Breaker(makeRequest, 3)

	for i := 0; i < 10; i++ {
		res, err := breaker(ctx)

		if err != nil {
			fmt.Printf("Error: '%s'\n", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		fmt.Println(res)
		break
	}

	res2, _ := breaker(ctx)
	fmt.Println(res2)
}
