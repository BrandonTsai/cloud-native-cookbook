package main

import (
	"fmt"
)

func main() {

	s := [][]int{
		{2, 4, 6, 8, 10},
		{1, 3, 5},
	}

	for i := 0; i < len(s); i++ {
		for _, value := range s[i] {
			fmt.Println(value)
		}
	}

}
