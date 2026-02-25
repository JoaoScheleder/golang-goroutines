package main

import "sync"

var msg string
var wg sync.WaitGroup

func updateMessage(s string) {
	defer wg.Done()
	msg = s
}

func updateMessageWithMutex(s string, mu *sync.Mutex) {
	defer wg.Done()
	mu.Lock()
	// msg is protected by the mutex, so only one goroutine can access it at a time
	msg = s
	mu.Unlock()
}

func RaceConditionFixedWithMutex() {
	// This function would use a mutex to protect access to the shared variable 'msg'
	var mu sync.Mutex
	msg = "Hello, World!"
	wg.Add(2)

	go updateMessageWithMutex("Hello, Go!", &mu)
	go updateMessageWithMutex("Hello, Concurrency!", &mu)

	wg.Wait()
	println(msg)
}

func RaceCondition() {
	msg = "Hello, World!"
	wg.Add(2)

	go updateMessage("Hello, Go!")
	go updateMessage("Hello, Concurrency!")

	wg.Wait()
	println(msg)

}
