package main

import (
	"errors"
	"math/rand"
)

type Deck struct {
	cards []string
}

func NewDeck(cards map[string]Card) *Deck {
	var d Deck
	for k, card := range cards {
		for i := 0; i < card.Count; i++ {
			d.Insert(k)
		}
	}
	return &d
}

func (d *Deck) Insert(card string) {
	d.cards = append(d.cards, card)
}

func (d *Deck) Shuffle(seed int64) {
	rand.Seed(seed)
	rand.Shuffle(len(d.cards), func(i, j int) { d.cards[i], d.cards[j] = d.cards[j], d.cards[i] })
}

func (d *Deck) Draw() (string, error) {
	if len(d.cards) < 1 {
		return "", errors.New("deck: no cards left")
	}
	card := d.cards[len(d.cards)-1]
	d.cards = d.cards[:len(d.cards)-1]
	return card, nil
}
