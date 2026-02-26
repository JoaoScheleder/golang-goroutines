package main

import (
	"fmt"
	"math/rand"
	"time"
)

// Global settings
const NumberOfPizzas = 10

// Global counters (tracked by the "kitchen")
var pizzasMade, pizzasFailed, total int

// Producer is like a manager holding two communication radios (channels)
type Producer struct {
	data chan PizzaOrder // A tube to send finished pizzas through
	quit chan chan error // A special tube used to signal "Time to close the shop!"
}

// PizzaOrder is the "receipt" for a single pizza
type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

// Close is a "clean-up" function to shut down the producer safely
func (p *Producer) Close() error {
	ch := make(chan error) // Create a temporary "feedback" tube
	p.quit <- ch           // Send that tube into the quit channel
	return <-ch            // Wait here until we get a response back
}

// Random number setup so every run feels different
var seed = time.Now().UnixNano()
var r = rand.New(rand.NewSource(seed))

// makePizza simulates the actual cooking process
func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++ // Move to the next order number

	if pizzaNumber <= NumberOfPizzas {
		delay := r.Intn(5) + 1 // Random cook time (1-5 seconds)
		fmt.Printf("Received order #%d!\n", pizzaNumber)

		rnd := r.Intn(12) + 1 // Random chance of success or failure
		msg := ""
		success := false

		// Phase 1: Check if we have ingredients
		if rnd < 5 {
			pizzasFailed++
			msg = fmt.Sprintf("Failed to make pizza #%d. Not enough ingredients.", pizzaNumber)
		} else {
			pizzasMade++
		}
		total++

		// Wait while the "pizza cooks"
		fmt.Printf("Making pizza #%d will take %d seconds\n", pizzaNumber, delay)
		time.Sleep(time.Duration(delay) * time.Second)

		// Phase 2: Check if the oven broke or staff left
		if rnd <= 2 {
			msg = fmt.Sprintf("&&& Failed to make pizza #%d. Oven is broken.", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("&&& Failed to make pizza #%d. Not enough staff.", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza #%d is ready!", pizzaNumber)
		}

		return &PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}
	}

	// If we've hit the limit (10), return a blank order
	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}
}

// pizzeria is the "Chef" routine that runs in the background
func pizzeria(pizzaMaker *Producer) {
	var i = 0

	for {
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber

			// 'select' is like a chef looking at two doors at once
			select {
			// DOOR 1: Put the pizza on the "delivery belt" (data channel)
			case pizzaMaker.data <- *currentPizza:

			// DOOR 2: If someone tells us to "Quit", clean up and go home
			case quitChan := <-pizzaMaker.quit:
				close(pizzaMaker.data) // Stop the belt
				close(quitChan)        // Confirm we are leaving
				return                 // Exit the loop/goroutine
			}
		}
	}
}

func main() {
	fmt.Println("The Pizzeria is open for business!")
	fmt.Println("----------------------------")

	// Create our "Manager" (Producer) with its communication tubes
	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	// 'go' starts a Goroutine. This is like hiring a chef
	// to work in the kitchen while you stay at the front desk.
	go pizzeria(pizzaJob)

	// Range over the channel: This is like a conveyor belt.
	// We stay here and grab every pizza that comes out.
	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				fmt.Printf("Order #%d: %s out for delivery!\n", i.pizzaNumber, i.message)
			} else {
				fmt.Printf("Order #%d: %s\n", i.pizzaNumber, i.message)
			}
		} else {
			// Once we reach 10, tell the kitchen to close
			fmt.Println("Done making pizzas")
			err := pizzaJob.Close()
			if err != nil {
				fmt.Printf("Error closing producer: %v\n", err)
			}
			break // Stop looking for more pizzas
		}
	}

	// Final Report
	fmt.Printf("Total pizzas made: %d\n", pizzasMade)
	fmt.Printf("Total pizzas failed: %d\n", pizzasFailed)
	fmt.Printf("Total orders processed: %d\n", total)
}
