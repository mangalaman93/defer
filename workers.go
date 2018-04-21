package main

import (
	"fmt"
	"sync"
	"time"
)

func doWork(w *sync.WaitGroup, id int) {
	defer w.Done()
	fmt.Println("Start worker", id)
	time.Sleep(time.Second / 2)
	fmt.Println("End worker", id)
}

func main() {
	var wg sync.WaitGroup

	wg.Add(2)
	go doWork(&wg, 0)
	go doWork(&wg, 1)

	fmt.Println("waiting for routines")
	wg.Wait()
	fmt.Println("exiting main")
}
