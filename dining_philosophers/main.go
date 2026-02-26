package main

import (
	"fmt"
	"sync"
	"time"
)

type Philosopher struct {
	name      string
	rightFork int
	leftFork  int
}

// list of all Philosophers
var philosophers = []Philosopher{
	{"Plato", 0, 1},
	{"Aristotle", 1, 2},
	{"Socrates", 2, 3},
	{"Confucius", 3, 4},
	{"Descartes", 4, 0},
}

var hunger = 3 // Each philosopher needs to eat 3 times
var eatTime = 1 * time.Second
var thinkTime = 3 * time.Second
var sleepTime = 1 * time.Second

func main() {
	fmt.Println("Dining Philosophers Problem Simulation")
	fmt.Println("-----------------------------------")
	fmt.Println("The table is empty")
	// start the meal
	dine()

	fmt.Println("The table is empty")

}

func dine() {
	wg := &sync.WaitGroup{}
	wg.Add(len(philosophers))

	seated := &sync.WaitGroup{}
	seated.Add(len(philosophers))

	// forks is a map of all 5 forks
	var forks = make(map[int]*sync.Mutex)
	for i := 0; i < len(philosophers); i++ {
		forks[i] = &sync.Mutex{}
	}

	// start the meal
	for i := 0; i < len(philosophers); i++ {
		// fire go routines for each philosopher
		go diningProblem(philosophers[i], forks, wg, seated)
	}
}

func diningProblem(philosopher Philosopher, forks map[int]*sync.Mutex, wg *sync.WaitGroup, seated *sync.WaitGroup) {
	defer wg.Done()

	// seat the philosopher at the table
	fmt.Printf("%s is seated at the table.\n", philosopher.name)
	seated.Done() // signal that this philosopher is seated
	seated.Wait() // wait for all philosophers to be seated
	// eat three times

	for i := 0; i < hunger; i++ {

		// fix logical race by picking up the right fork first, then the left fork
		if philosopher.leftFork < philosopher.rightFork {

			forks[philosopher.rightFork].Lock() // pick up right fork
			fmt.Printf("%s picked up right fork %d.\n", philosopher.name, philosopher.rightFork)

			forks[philosopher.leftFork].Lock() // pick up left fork
			fmt.Printf("%s picked up left fork %d.\n", philosopher.name, philosopher.leftFork)
		} else {
			forks[philosopher.leftFork].Lock()
			fmt.Printf("%s picked up left fork %d.\n", philosopher.name, philosopher.leftFork)

			forks[philosopher.rightFork].Lock() // pick up right fork
			fmt.Printf("%s picked up right fork %d.\n", philosopher.name, philosopher.rightFork)
		}
		// Original race condition
		// forks[philosopher.rightFork].Lock() // pick up right fork
		// fmt.Printf("%s picked up right fork %d.\n", philosopher.name, philosopher.rightFork)

		// forks[philosopher.leftFork].Lock()
		// fmt.Printf("%s picked up left fork %d.\n", philosopher.name, philosopher.leftFork)

		fmt.Printf("%s is eating.\n", philosopher.name)
		time.Sleep(eatTime) // simulate eating

		// thinking
		fmt.Printf("%s is thinking.\n", philosopher.name)
		time.Sleep(thinkTime) // simulate thinking

		// put down the forks
		forks[philosopher.rightFork].Unlock()
		fmt.Printf("%s put down right fork %d.\n", philosopher.name, philosopher.rightFork)

		forks[philosopher.leftFork].Unlock()
		fmt.Printf("%s put down left fork %d.\n", philosopher.name, philosopher.leftFork)
	}

	// philosopher is satisfied and leaves the table
	fmt.Printf("%s is satisfied and leaves the table.\n", philosopher.name)
}
