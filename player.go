package main

import (
	"errors"
	"fmt"
)

type Player struct {
	Id     int
	Name   string
	Status Status

	Distance   int
	Score      int
	RoundBonus int
	Immunities []Status

	hand    []Card
	terrain []Card
	battle  []Card
	defense []Card
	travel  []Card
}

// Play plays the selected card from player hand into field
// TODO: move hand logic to separate object
func (p *Player) Play(c int) (Card, error) {
	if c > len(p.hand)-1 {
		return Card{}, errors.New("hand: invalid card index")
	}
	cardPlayed := p.hand[c]
	p.hand = append(p.hand[:c], p.hand[c+1:]...)
	return cardPlayed, nil
}

// Receive receives a card from self (hand) or another player and put on the player field
func (p *Player) Receive(from *Player, card Card) error {
	// validate constraints
	// constraint type: target
	if self := p == from; card.Constraints.Target != "" && (self && card.Constraints.Target != "self" || !self && card.Constraints.Target == "self") {
		return fmt.Errorf("can't play card %s on %s field: invalid target", card.Name, p.Name)
	}

	// constraint type: status
	if len(card.Constraints.Status) > 0 {
		validStatus := false
		for _, s := range card.Constraints.Status {
			if s == p.Status {
				validStatus = true
				break
			}
		}

		if !validStatus {
			return fmt.Errorf("can't play card %s on status %s", card.Name, p.Status)
		}
	}

	// constraint type: terrain
	if len(card.Constraints.Terrains) > 0 && len(p.terrain) > 0 {
		validTerrain := false
		for _, t := range card.Constraints.Terrains {
			if t == p.Terrain().Effects.Terrain {
				validTerrain = true
				break
			}
		}

		if !validTerrain {
			return fmt.Errorf("can't play card %s on terrain %s", card.Name, p.Terrain().Effects.Terrain)
		}
	}

	// validate immunities
	if len(p.Immunities) > 0 {
		for _, i := range p.Immunities {
			immune := false
			if card.Effects.Status == i || (i == StatusLost && card.Type == TypeYellow) {
				immune = true
				break
			}

			if immune {
				return fmt.Errorf("can't play card %s, player has %s immunity", card.Name, i)
			}
		}
	}

	// add card to proper stack
	switch card.Type {
	case TypeBlue:
		p.defense = append(p.defense, card)
	case TypeGreen, TypeRed:
		p.battle = append(p.battle, card)
	case TypeWhite:
		p.travel = append(p.travel, card)
	case TypeYellow:
		p.terrain = append(p.terrain, card)
	}

	// calculate effects
	p.Distance += card.Effects.Distance
	p.RoundBonus += card.Effects.Bonus
	if card.counter {
		p.RoundBonus += card.Effects.CounterBonus
	}

	// blue card effects
	if card.Effects.Immunity != StatusUnknown {
		p.Immunities = append(p.Immunities, card.Effects.Immunity)

		if p.Status == card.Effects.Immunity {
			// remove current status
			// TODO: add removed card to graveyard
			p.battle = p.battle[:len(p.battle)-1]
		}
	}

	// status
	if len(p.battle) > 0 {
		p.Status = p.Battle().Effects.Status
	}

	// apply RA special effects
	for _, i := range p.Immunities {
		// immune to status lost == always oriented
		if i == StatusLost {
			switch p.Status {
			case StatusEscaping, StatusHealing, StatusWorking:
				p.Status = StatusOriented
			}

			// Remove all terrains
			// TODO: move them to graveyard
			p.terrain = nil
		}
	}

	return nil
}

// Hand returns the cards in the player hand
func (p *Player) Hand() []Card {
	return p.hand
}

// Defense returns the defense cards in the player field
func (p *Player) Defense() []Card {
	return p.defense
}

// Draw draws a card from the deck and add to the players hand
func (p *Player) Draw(d *Deck) error {
	card, err := d.Draw()
	if err != nil {
		return err
	}
	p.hand = append(p.hand, card)
	return nil
}

// Terrain returns the top card of the terrain stack
func (p *Player) Terrain() Card {
	if len(p.terrain) == 0 {
		return Card{}
	}
	return p.terrain[len(p.terrain)-1]
}

// Battle returns the top card of the battle stack
func (p *Player) Battle() Card {
	if len(p.battle) == 0 {
		return Card{}
	}
	return p.battle[len(p.battle)-1]
}
