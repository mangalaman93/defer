package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	errFailedBilling      = errors.New("error in billing customer")
	errFailedNotification = errors.New("error in notifying customer")

	randSrc = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type customer struct {
	mutex   sync.Mutex
	balance int
}

func (c *customer) Bill(amount int) error {
	// if randSrc.Intn(100) < 50 {
	// 	return errFailedBilling
	// }

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.balance -= amount

	fmt.Printf("customer has been billed $%v\n", amount)
	return nil
}

func (c *customer) Notify() error {
	// if randSrc.Intn(100) < 50 {
	// 	return errFailedNotification
	// }

	fmt.Println("customer has been notified")
	return nil
}

func billCustomer(c *customer, amount int) {
	if err := c.Bill(amount); err != nil {
		fmt.Println("unable to bill customer ::", err)
		return
	}

	if err := c.Notify(); err != nil {
		fmt.Println("unable to notify customer, saving notification ::", err)
		return
	}

	fmt.Println("billed and notified customer")
}

func main() {
	customer := &customer{balance: 1000}
	fmt.Println("initial balance: $", customer.balance)

	go billCustomer(customer, 100)
	go billCustomer(customer, 200)
	go billCustomer(customer, 400)
	time.Sleep(1 * time.Second)
	fmt.Println("final balance: $", customer.balance)
}
