package main

import "github.com/hajimehoshi/ebiten/v2"

type Hand struct {
	cards    []Card
	op       *ebiten.DrawImageOptions
	selected int
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

// Sprite returns the next frame in animation
func (h *Hand) Sprites() ([]*ebiten.Image, []*ebiten.DrawImageOptions) {
	return nil, nil
}
