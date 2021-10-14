package main

import (
	"errors"
	"log"
)

type Player struct {
	hand    []Card
	terrain []Card
	battle  []Card
	defense []Card
	travel  []Card
	score   int
}

func (p *Player) Hand() []Card {
	return p.hand
}

func (p *Player) Draw(d *Deck) error {
	card, err := d.Draw()
	log.Printf("card drawn: %v", card)
	if err != nil {
		return err
	}
	p.hand = append(p.hand, card)
	return nil
}

func (p *Player) Terrain() Card {
	if len(p.terrain) == 0 {
		return Card{}
	}
	return p.terrain[len(p.terrain)-1]
}

func (p *Player) BattleStatus() string {
	if len(p.battle) == 0 {
		for _, c := range p.defense {
			if c.Key == "RA" {
				return "orientation"
			}
		}
		return "lost"
	}
	return p.battle[len(p.battle)-1].Effect
}

func (p *Player) Battle() Card {
	if len(p.battle) == 0 {
		return Card{}
	}
	return p.battle[len(p.battle)-1]
}

func (p *Player) Play(card int) (Card, error) {
	if card > len(p.hand)-1 {
		return Card{}, errors.New("hand: invalid card index")
	}
	cardPlayed := p.hand[card]
	p.hand = append(p.hand[:card], p.hand[card+1:]...)
	return cardPlayed, nil
}

func (p *Player) Receive(card Card) error {
	switch card.Type {
	case "green":
		p.battle = append(p.battle, card)
	case "yellow":
		p.terrain = append(p.terrain, card)
	}
	return nil
}
