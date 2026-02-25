package main

import (
	"fmt"
	"sync"
)

type Income struct {
	Source string
	Amount int
}

var wg sync.WaitGroup

func main() {
	var bankBalance int = 0
	var balance sync.Mutex
	incomes := []Income{
		{Source: "Job", Amount: 5000},
		{Source: "Freelance", Amount: 2000},
		{Source: "Investments", Amount: 1000},
	}

	wg.Add(len(incomes))
	for i, income := range incomes {
		go func(i int, income Income) {
			defer wg.Done()
			for week := 1; week <= 52; week++ {
				balance.Lock()
				bankBalance += income.Amount
				fmt.Printf("Week %d: %s - $%d\n", week, income.Source, income.Amount)
				balance.Unlock()
			}
		}(i, income)
	}

	wg.Wait()

	fmt.Printf("Final Bank Balance: $%d\n", bankBalance)
}
