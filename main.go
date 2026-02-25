package main

import (
	"fmt"
	"time"
)

func printSomething(s string) {
	fmt.Println(s)
}

func main() {

	words := []string{"Hello", "World", "Go", "Programming"}

	for i, word := range words {
		go printSomething(fmt.Sprintf("%d: %s", i, word))
	}

	time.Sleep(1 * time.Second)
}
