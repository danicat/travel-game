package main

import "github.com/hajimehoshi/ebiten/v2"

type Stack struct {
	cards []Card
	op    *ebiten.DrawImageOptions
}

func (s *Stack) Put(c Card) {
	s.cards = append(s.cards, c)
}

func (s *Stack) Remove() Card {
	if len(s.cards) == 0 {
		return Card{}
	}
	card := s.Top()
	s.cards = s.cards[:len(s.cards)-1]
	return card
}

func (s *Stack) Top() Card {
	if len(s.cards) == 0 {
		return Card{}
	}
	card := s.cards[len(s.cards)-1]
	return card
}

func (s *Stack) All() []Card {
	return s.cards
}

// Sprite returns the next frame in animation
func (s *Stack) Sprite() (*ebiten.Image, *ebiten.DrawImageOptions) {
	return s.Top().image, s.op
}

type BattleStack struct {
	Stack
}

func (bs *BattleStack) Status() Status {
	return bs.Top().Effects.Status
}

type TerrainStack struct {
	Stack
}

func (ts *TerrainStack) Terrain() string {
	return ts.Top().Effects.Terrain
}

func (ts *TerrainStack) Clear() {
	ts.cards = nil
}

type DefenseStack struct {
	Stack
}

func (ds *DefenseStack) Immunities() []Status {
	var im []Status
	for _, c := range ds.cards {
		im = append(im, c.Effects.Immunity)
	}
	return im
}
