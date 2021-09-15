package main

import (
	"fmt"
)

func main() {
	hello()

	x := 1
	y := []int{5, 4, 3, 2, 1}
	reset(x, y)
	fmt.Println(x, y)
}

func hello() {
	fmt.Println("Hello world")
}

func reset(i int, a []int) {
	i = 0
	a[0] = 1
}
