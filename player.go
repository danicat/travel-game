package main

import (
	"errors"
	"fmt"
	"log"
)

type Player struct {
	Id         int
	Name       string
	Status     PlayerStatus
	Score      int
	RoundBonus int

	hand    []Card
	terrain []Card
	battle  []Card
	defense []Card
	travel  []Card
}

func (p *Player) Hand() []Card {
	return p.hand
}

func (p *Player) Draw(d *Deck) error {
	card, err := d.Draw()
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

// Receive receives a card from self (hand) or another player and put on the player field
func (p *Player) Receive(from *Player, card Card) error {
	self := p == from
	switch card.Type {
	case TypeBlue:
		if !self {
			return errors.New("can't play blue cards on opponent's field")
		}
		p.RoundBonus += 4000
		p.defense = append(p.defense, card)

	case TypeWhite:
		if !self {
			return errors.New("can't play white cards on opponent's field")
		}
		validMove := false
		// check constraints

		if !validMove {
			return fmt.Errorf("can't play %s on %s terrain", card.Name, p.Terrain().Name)
		}

		p.Score += card.Distance
		p.travel = append(p.travel, card)

	case TypeGreen:
		if !self {
			return errors.New("can't play green cards on opponent's field")
		}
		validMove := false
		// check constraints

		if !validMove {
			return fmt.Errorf("can't play %s on %s status", card.Name, p.Status)
		}

		p.battle = append(p.battle, card)

	case TypeRed:
		if !self {
			return errors.New("can't play red cards on self's field")
		}

		validMove := false
		// check constraints

		if !validMove {
			return fmt.Errorf("can't play %s on %s status", card.Name, p.Status)
		}

		p.battle = append(p.battle, card)

	case TypeYellow:
		p.terrain = append(p.terrain, card)
	}
	return nil
}

type Players struct {
	players       []*Player
	currentPlayer int
}

func NewPlayers(numPlayers int) *Players {
	if numPlayers == 0 {
		log.Fatal("cannot start game with zero players")
	}

	var players []*Player
	for i := 0; i < numPlayers; i++ {
		name := fmt.Sprintf("Player %d", i+1)
		players = append(players, &Player{Id: i, Name: name})
	}
	return &Players{players: players}
}

func (p *Players) Current() *Player {
	return p.players[p.currentPlayer]
}

// Next returns the next player in the turn cycle and advance the current player index
func (p *Players) Next() *Player {
	if p.currentPlayer == len(p.players)-1 {
		p.currentPlayer = 0
	} else {
		p.currentPlayer++
	}
	return p.players[p.currentPlayer]
}

// PeekNext returns the next player in the turn cycle but doesn't advance the current player index
func (p *Players) PeekNext() *Player {
	if p.currentPlayer == len(p.players)-1 {
		return p.players[0]
	}
	return p.players[p.currentPlayer+1]
}

func (p *Players) Len() int {
	return len(p.players)
}

func (p *Players) All() []*Player {
	return p.players
}
