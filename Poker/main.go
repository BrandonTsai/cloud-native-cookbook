package main

import (
	"fmt"

	"brandon/poker/deck"
)

func main() {
	a := deck.GetStandardDeck()
	for _, x := range a {
		fmt.Println(x)
	}
}
