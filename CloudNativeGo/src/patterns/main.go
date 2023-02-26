package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
)

func makeRequest(ctx context.Context) (string, error) {
	resp, err := http.Get("http://localhost:8080/albums")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	s := string(body)
	return s, nil
}

func main() {
	ctx := context.TODO()
	breaker := Breaker(makeRequest, 3)
	res, err := breaker(ctx)

	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
