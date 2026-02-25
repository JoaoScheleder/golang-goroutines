package main

import (
	"fmt"
	"math/rand"
	"time"
)

const NumberOfPizzas = 10

var pizzasMade, pizzasFailed, total int

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

func (p *Producer) Close() error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

var seed = time.Now().UnixNano()
var r = rand.New(rand.NewSource(seed))

func makePizza(pizzaNumber int) *PizzaOrder {
	pizzaNumber++
	if pizzaNumber <= NumberOfPizzas {
		delay := r.Intn(5) + 1
		fmt.Printf("Receveid order #%d!\n", pizzaNumber)
		rnd := r.Intn(12) + 1
		msg := ""
		success := false

		if rnd < 5 {
			pizzasFailed++
			msg = fmt.Sprintf("Failed to make pizza #%d. Not enough ingredients.", pizzaNumber)
		} else {
			pizzasMade++
		}
		total++

		fmt.Printf("Making pizza #%d will take %d seconds\n", pizzaNumber, delay)
		time.Sleep(time.Duration(delay) * time.Second)

		if rnd <= 2 {
			msg = fmt.Sprintf("&&& Failed to make pizza #%d. Oven is broken.", pizzaNumber)
		} else if rnd <= 4 {
			msg = fmt.Sprintf("&&& Failed to make pizza #%d. Not enough staff.", pizzaNumber)
		} else {
			success = true
			msg = fmt.Sprintf("Pizza #%d is ready!", pizzaNumber)
		}

		p := PizzaOrder{
			pizzaNumber: pizzaNumber,
			message:     msg,
			success:     success,
		}

		return &p
	}

	return &PizzaOrder{
		pizzaNumber: pizzaNumber,
	}

}

func pizzeria(pizzaMaker *Producer) {
	var i = 0

	for {
		currentPizza := makePizza(i)
		if currentPizza != nil {
			i = currentPizza.pizzaNumber
			select {
			case pizzaMaker.data <- *currentPizza:

			case quitChan := <-pizzaMaker.quit:
				close(pizzaMaker.data)
				close(quitChan)
				return
			}
		}
	}
}

func main() {

	fmt.Println("The Pizzeria is open to business!")
	fmt.Println("----------------------------")

	pizzaJob := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	go pizzeria(pizzaJob)

	for i := range pizzaJob.data {
		if i.pizzaNumber <= NumberOfPizzas {
			if i.success {
				fmt.Printf("Order #%d: %s out for delivery!\n", i.pizzaNumber, i.message)
			} else {
				fmt.Printf("Order #%d: %s\n failed", i.pizzaNumber, i.message)
			}
		} else {
			fmt.Println("Done making pizzas")
			err := pizzaJob.Close()
			if err != nil {
				fmt.Printf("Error closing producer: %v\n", err)
			}

			break
		}
	}

	fmt.Printf("Total pizzas made: %d\n", pizzasMade)
	fmt.Printf("Total pizzas failed: %d\n", pizzasFailed)
	fmt.Printf("Total orders processed: %d\n", total)

}
