package main

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

var (
	ErrFailedBilling      = errors.New("error in billing customer")
	ErrFailedNotification = errors.New("error in notifying customer")

	randSrc = rand.New(rand.NewSource(time.Now().UnixNano()))
)

type Customer struct {
	mutex   sync.Mutex
	balance int
}

func (c *Customer) Bill(amount int) error {
	//if randSrc.Intn(100) < 50 {
	//	return ErrFailedBilling
	//}

	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.balance -= amount

	fmt.Printf("customer has been billed $%v\n", amount)
	return nil
}

func (c *Customer) Notify() error {
	//if randSrc.Intn(100) < 50 {
	//	return ErrFailedNotification
	//}

	fmt.Println("customer has been notified")
	return nil
}

func BillCustomer(c *Customer, amount int) {
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
	customer := &Customer{balance: 1000}
	go BillCustomer(customer, 100)
	go BillCustomer(customer, 200)
	go BillCustomer(customer, 400)
	time.Sleep(1 * time.Second)
	fmt.Println("final balance:", customer.balance)
}
