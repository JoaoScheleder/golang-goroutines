package main

import (
	"fmt"
	"math/rand"
	"time"
)

var seatingCapacity = 10
var arrivalRate = 1000
var cutDuration = 1000 * time.Millisecond
var timeOpen = 8 * time.Second

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

type BarberShop struct {
	ShopCapacity    int
	HairCutDuration time.Duration
	NumberOfBarbers int
	BarbersDoneChan chan bool
	ClientsChan     chan string
	Open            bool
}

func (b *BarberShop) addBarber(barber string) {
	b.NumberOfBarbers++

	go func() {
		isSleeping := false
		fmt.Printf("Barber %d is ready to cut hair, and are in the room to check for clients	\n", b.NumberOfBarbers)
		for {
			// if there are no clients, the barber goes to sleep
			if len(b.ClientsChan) == 0 {
				fmt.Printf("Barber %d is sleeping\n", b.NumberOfBarbers)
				isSleeping = true
			}

			client, shopOpen := <-b.ClientsChan
			if shopOpen {
				if isSleeping {
					fmt.Printf("Barber %d is waking up to cut client %d hair\n", b.NumberOfBarbers, client)
					isSleeping = false
				}
				b.cutHair(barber, fmt.Sprintf("Client %d", client))
			} else {
				b.sendBarberHome()

			}
		}
	}()
}

func (b *BarberShop) cutHair(barber, client string) {
	fmt.Printf("Barber %s is cutting hair of client %s\n", barber, client)
	time.Sleep(b.HairCutDuration)
	fmt.Printf("Barber %s is done cutting hair of client %s\n", barber, client)
}

func (b *BarberShop) sendBarberHome() {
	fmt.Printf("Barber %d is going home\n", b.NumberOfBarbers)
	b.BarbersDoneChan <- true
}

func (b *BarberShop) closeShopForTheDay() {
	fmt.Println("Shop is closing for the day")
	b.Open = false
	close(b.ClientsChan)
	for i := 0; i < b.NumberOfBarbers; i++ {
		<-b.BarbersDoneChan
	}
	fmt.Println("---------------------------------------")
	fmt.Println("Shop is closed for the day")
}

func (b *BarberShop) addClient(client string) {
	fmt.Println("Client arrived %d", client)
	if b.Open {
		select {
		case b.ClientsChan <- client:
			fmt.Printf("Client %d take a seat in the waiting room\n", client)
		default:
			fmt.Printf("Shop is full, client %s is leaving\n", client)
		}
	} else {
		fmt.Printf("Shop is closed, client %s is leaving\n", client)
	}
}

func main() {
	fmt.Println("Sleeping barber problem")

	clientChan := make(chan string, seatingCapacity)
	doneChan := make(chan bool)

	shop := BarberShop{
		ShopCapacity:    seatingCapacity,
		HairCutDuration: cutDuration,
		NumberOfBarbers: 0,
		BarbersDoneChan: doneChan,
		ClientsChan:     clientChan,
		Open:            true,
	}

	fmt.Println("Shop is open")

	// add barbers
	shop.addBarber("Frank")
	shop.addBarber("Luis")
	shop.addBarber("Marcus")
	shop.addBarber("Matheus")
	shop.addBarber("Leonardo")
	shop.addBarber("Aurelio")
	shop.addBarber("Anderson")
	shop.addBarber("Luciano")
	shop.addBarber("Filipe")

	shopClosing := make(chan bool)
	closed := make(chan bool)

	go func() {
		<-time.After(timeOpen)
		shopClosing <- true
		closed <- true
	}()

	i := 1

	go func() {
		for {
			randomMilis := r.Intn(10000) % arrivalRate
			select {
			case <-shopClosing:
				return

			case <-time.After(time.Duration(randomMilis) * time.Millisecond):
				shop.addClient(fmt.Sprintf("Client %d", i))
				i++
			}
		}
	}()

	time.Sleep(5 * time.Second)

}
