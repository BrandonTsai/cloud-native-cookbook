package main

import "fmt"

func main() {

	b := map[string]int{
		"Bubble Milk Tea": 55,
		"Ice Coffee":      45,
	}
	printMap(b)
}

func printMap(p map[string]int) {
	for key, value := range p {
		fmt.Printf("The price of %s is %d\n", key, value)
	}
}
