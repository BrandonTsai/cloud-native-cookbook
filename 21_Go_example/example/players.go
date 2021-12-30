package main

type contactInfo struct {
	phone string
	email string
}

type player struct {
	name    string
	coins   int
	contact contactInfo
}

func (a player) isRicherThan(b player) bool {
	if a.coins > b.coins {
		return true
	}
	return false
}

func (p1 *player) giveCoinsTo(p2 *player, payCoins int) {
	(*p1).coins = (*p1).coins - payCoins
	(*p2).coins = (*p2).coins + payCoins
}
