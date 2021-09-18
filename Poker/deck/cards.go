package deck

func GetStandardDeck() [52]string {
	cards := [52]string{}

	i := 0
	for _, s := range getSuits() {
		for _, r := range getRanks() {
			cards[i] = s + " " + r
			i += 1
		}
	}

	return cards
}

func getSuits() [4]string {
	return [4]string{"Clubs", "Diamonds", "Hearts", "Spades"}
}

func getRanks() [13]string {
	return [13]string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K"}
}
