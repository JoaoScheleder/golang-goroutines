package main

import (
	"fmt"
	"time"
)

func listenToChannel(ch chan string) {
	for {
		i := <-ch
		fmt.Println("Got", i, " from channel")

		// simulate doing a lot of work
		time.Sleep(1 * time.Second)
	}
}

func main() {
	// ch := make(chan string, 100)
	ch := make(chan string)
	go listenToChannel(ch)

	// buffered channel can hold 100 messages, so we can send 100 messages without blocking
	// unbuffered channel will block until the listener receives the message, so we can only send one message at a time

	for i := 0; i < 100; i++ {
		fmt.Println("Sending", i, " To channel...")
		ch <- fmt.Sprintf("Message %d", i)
		fmt.Println("Sent", i, " To channel")
	}
}
