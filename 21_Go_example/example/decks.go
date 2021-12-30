package main

import "fmt"

type card struct {
	suit string
	rank string
}

// Create a new type of 'deck'
// which is a slice of strings
type decks []card

// 'd' is the copy of the variable of type 'decl'
// that has access to this NewStandardDeck() method
func (d decks) newStandardDeck() decks {
	newDecks := []card{}
	suits := [4]string{"Clubs", "Diamonds", "Hearts", "Spades"}
	ranks := [13]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}

	for _, s := range suits {
		for _, r := range ranks {

			// you can initial card via "card{suit: s, rank: r}" as well
			newDecks = append(newDecks, card{s, r})
		}
	}
	return newDecks
}

func (d decks) print() {
	for _, value := range d {
		fmt.Println(value.suit, value.rank)
	}
}

func (d decks) printWithField() {
	for _, value := range d {
		fmt.Printf("%+v\n", value)
	}
}
