package main

import (
	"fmt"
	"time"
)

func main() {
	messages := make(chan string)
	signals := make(chan bool)

	// go func() {
	// 	messages <- "message to send"
	// 	signals <- true
	// }()

	time.Sleep(100 * time.Millisecond)

	select {
	case msg := <-messages:
		fmt.Printf("received message: \"%s\"\n", msg)
	default:
		fmt.Println("no messages received")
	}

	go func() {
		msg := "hi"
		select {
		case messages <- msg:
			fmt.Println("sent message:", msg)
		default:
			fmt.Println("no message sent")
		}
	}()

	//fmt.Println("received:", <-messages)
	time.Sleep(100 * time.Millisecond)

	select {
	case msg := <-messages:
		fmt.Println("received message:", msg)
	case sig := <-signals:
		fmt.Println("received signal:", sig)
	default:
		fmt.Println("no activity")
	}
}
