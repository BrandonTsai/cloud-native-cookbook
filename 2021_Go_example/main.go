package main

import "fmt"

func main() {
	for i := 100; i < 150; i++ {
		if i%2 == 0 {
			continue
		} else if i%7 == 0 {
			fmt.Printf("%d is a the first multiple of 7\n", i)
			break
		}
		fmt.Printf("%d is a the odd number\n", i)
	}
}
