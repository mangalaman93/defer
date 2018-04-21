package main

import (
	"fmt"
	"log"
	"sync/atomic"
	"time"
)

const (
	cStatusWaiting    = 0
	cStatusInProgress = 1
)

var (
	status int32
)

func doCrucialWork() {
	log.Println("[INFO] starting work")
	time.Sleep(time.Second * 2)
	log.Println("[INFO] work complete")
}

func startWork() {
	if atomic.SwapInt32(&status, cStatusInProgress) == cStatusInProgress {
		log.Println("[INFO] work is already in progress")
		return
	}

	defer atomic.SwapInt32(&status, cStatusWaiting)
	doCrucialWork()
}

func printStatus() {
	curStatus := atomic.LoadInt32(&status)
	if curStatus == cStatusWaiting {
		fmt.Println("Current Status: waiting")
	} else if curStatus == cStatusInProgress {
		fmt.Println("Current Status: inProgress")
	} else {
		panic("SHOULD NOT REACH")
	}
}

func main() {
	count := 0
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			printStatus()
			count = count + 1
			log.Println("[INFO] received ticker to do work")
			go startWork()
			if count == 10 {
				return
			}
		}
	}

	time.Sleep(time.Second * 2)
}
