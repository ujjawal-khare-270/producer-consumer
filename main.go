package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Producer struct {
	data chan PizzaOrder
	quit chan chan error
}

type PizzaOrder struct {
	pizzaNumber int
	message     string
	success     bool
}

const numMaxSuccessfulPizza = 13

var totalSuccess int

func preparePizza(pizzaNumber int) PizzaOrder {
	randomNumber := rand.Intn(10) + 1
	pizzaNumber++

	time.Sleep(time.Second*time.Duration(randomNumber) + time.Millisecond*time.Duration(randomNumber))
	if randomNumber < 7 {
		fmt.Println("random number ", randomNumber)
		totalSuccess = totalSuccess + 1
		fmt.Sprintf("Pizza %v success \n", pizzaNumber)
		return PizzaOrder{pizzaNumber: pizzaNumber, message: "Success", success: true}
	} else {
		fmt.Sprintf("Pizza %v failed \n", pizzaNumber)
		return PizzaOrder{pizzaNumber: pizzaNumber, message: "Failure", success: false}
	}
}

func makePizza(p *Producer) {
	i := 0
	for {
		pizzaOrder := preparePizza(i)
		i = pizzaOrder.pizzaNumber
		if pizzaOrder.success {
			select {
			case p.data <- pizzaOrder:

			case quitChan := <-p.quit:
				close(p.data)
				close(quitChan)
			}
		}
	}
}

func Close(p *Producer) error {
	ch := make(chan error)
	p.quit <- ch
	return <-ch
}

func main() {
	rand.Seed(time.Now().UnixNano())
	pizzaProducer := &Producer{
		data: make(chan PizzaOrder),
		quit: make(chan chan error),
	}

	go makePizza(pizzaProducer)

	for i := range pizzaProducer.data {
		if totalSuccess > numMaxSuccessfulPizza {
			fmt.Println("Pizza making completed")
			Close(pizzaProducer)
		} else {
			fmt.Printf("%v pizza completed %d \n", totalSuccess, i.pizzaNumber)
		}
	}
}
