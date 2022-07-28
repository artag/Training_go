package main

import "fmt"

func main() {
	messages := make(chan string)

	go func() { messages <- "pong" }()

	fmt.Println("ping")

	msg := <-messages
	fmt.Println(msg)
}
