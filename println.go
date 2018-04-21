package main

import (
	"fmt"
)

func main() {
	fmt.Println("before defer")
	defer fmt.Println("in defer")
	fmt.Println("after defer")
}
