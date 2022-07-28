package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	c1 := make(chan string)
	c2 := make(chan string)

	go func() {
		sendMessageToChan("func1", "one", c1)
	}()

	go func() {
		sendMessageToChan("func2", "two", c2)
	}()

	for i := 0; i < 2; i++ {
		select {
		case msg1 := <-c1:
			fmt.Println("received:", msg1)
		case msg2 := <-c2:
			fmt.Println("received:", msg2)
		}
	}
}

func sendMessageToChan(info string, msg string, ch chan string) {
	rnd := rand.Intn(5)
	fmt.Println("func", info, "duration:", rnd, "seconds")
	k := time.Duration(rnd)
	time.Sleep(k * time.Second)
	ch <- msg
}
