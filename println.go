package main

import (
	"fmt"
)

func main() {
	fmt.Println("before defer")
	defer fmt.Println("in defer")
	defer fmt.Println("in defer second one")
	fmt.Println("after defer")
}
