package main

import (
	"fmt"
)

func main() {

	s := []int{2, 4, 6, 8, 10}

	// get value with index
	for i := 0; i < len(s); i++ {
		fmt.Println(s[i])
	}

	// you can use 'range' operator as well.
	for _, value := range s {
		fmt.Println(value)
	}

}
