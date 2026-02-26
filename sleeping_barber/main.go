package main

import (
	"fmt"
	"strings"
)

// ping send only channel
// pong receive only channel
func shout(ping <-chan string, pong chan<- string) {
	for {
		s := <-ping
		pong <- fmt.Sprintf("%s!!!", strings.ToUpper(s))
	}
}

func main() {
	ping := make(chan string)
	pong := make(chan string)

	go shout(ping, pong)

	fmt.Println("Type something and Press enter:")

	for {
		fmt.Printf("->")
		var input string
		fmt.Scanln(&input)

		if strings.ToLower(input) == "q" {
			break
		}

		ping <- input

		response := <-pong
		fmt.Printf("Response: %s\n", response)
	}

	fmt.Println("All done, closing channels")
	close(ping)
	close(pong)
}
