package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Hand struct {
	cards    []Card
	selected int
}

func NewHand() *Hand {
	return &Hand{}
}

func (h *Hand) All() []Card {
	return h.cards
}

func (h *Hand) Put(c Card) {
	h.cards = append(h.cards, c)
}

func (h *Hand) Remove() *Card {
	if len(h.cards) == 0 {
		return nil
	}
	card := h.Selected()
	h.cards = append(h.cards[:h.selected], h.cards[h.selected+1:]...)
	if h.selected > len(h.cards)-1 {
		h.selected = len(h.cards) - 1
	}
	return card
}

func (h *Hand) Selected() *Card {
	if len(h.cards) == 0 {
		return nil
	}
	card := h.cards[h.selected]
	return &card
}

func (h *Hand) Left() {
	if h.selected <= 0 {
		h.selected = len(h.cards) - 1
		return
	}

	h.selected--
}

func (h *Hand) Right() {
	if h.selected >= len(h.cards)-1 {
		h.selected = 0
		return
	}

	h.selected++
}

func (h *Hand) Draw(target *ebiten.Image) {
	// angle := (180 / math.Pi) / float64(len(h.cards))
	for i, c := range h.cards {
		var scale float64
		if h.selected == i {
			scale = .12
		} else {
			scale = .10
		}

		op := ebiten.DrawImageOptions{}
		op.GeoM.Scale(scale, scale)
		// op.GeoM.Rotate(float64(-3*i) * angle)
		op.GeoM.Translate(config.Layout.Hand.StartX, config.Layout.Hand.StartY)
		op.GeoM.Translate(float64(i)*config.Layout.Card.Width, 0)
		target.DrawImage(c.image, &op)
	}
}
