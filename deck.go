package main

import (
	"errors"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Deck struct {
	cards    []Card
	cardBack *ebiten.Image
	op       *ebiten.DrawImageOptions
}

func NewDeck(cards []Card) (*Deck, error) {
	var d Deck
	for _, card := range cards {
		for i := 0; i < card.Count; i++ {
			d.Insert(card)
		}
	}

	img, _, err := ebitenutil.NewImageFromFile("assets/cards/back.png")
	if err != nil {
		return nil, err
	}

	d.cardBack = img
	d.op = &ebiten.DrawImageOptions{}
	d.Reset()

	return &d, nil
}

func (d *Deck) Insert(card Card) {
	d.cards = append(d.cards, card)
}

func (d *Deck) Shuffle(seed int64) {
	rand.Seed(seed)
	rand.Shuffle(len(d.cards), func(i, j int) { d.cards[i], d.cards[j] = d.cards[j], d.cards[i] })
}

func (d *Deck) DrawCard() (Card, error) {
	if len(d.cards) < 1 {
		return Card{}, errors.New("deck: no cards left")
	}
	card := d.cards[len(d.cards)-1]
	d.cards = d.cards[:len(d.cards)-1]
	return card, nil
}

func (d *Deck) Reset() {
	d.op.GeoM.Reset()
	d.op.GeoM.Scale(.125, .15)
	d.op.ColorM.ChangeHSV(0, 1, .8)
	d.op.GeoM.Translate(config.Layout.System.StartX+config.Layout.Deck.StartX, config.Layout.System.StartY+config.Layout.Deck.StartY)
}

func (d *Deck) Draw(target *ebiten.Image) {
	if len(d.cards) > 0 {
		target.DrawImage(d.cardBack, d.op)
	}
}
