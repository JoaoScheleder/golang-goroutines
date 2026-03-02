package main

import (
	"fmt"
	"time"
)

func server1(ch chan<- string) {
	for {
		time.Sleep(10 * time.Second)
		ch <- "This is from server 1"
	}
}

func server2(ch chan<- string) {
	for {
		time.Sleep(5 * time.Second)
		ch <- "This is from server 2"
	}
}

func main() {
	fmt.Println("Select with channels")
	fmt.Println("---------------------")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go server1(ch1)
	go server2(ch2)

	for {
		select {
		case msg1 := <-ch1:
			fmt.Printf("1 Received: %s\n", msg1)
		case msg2 := <-ch2:
			fmt.Printf("2 Received: %s\n", msg2)
		case msg4 := <-ch1:
			fmt.Printf("3 Received: %s\n", msg4)
		case msg5 := <-ch2:
			fmt.Printf("4 Received: %s\n", msg5)
		}
	}
}
