package main

import (
	"errors"
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
)

type Player struct {
	ID         int
	Name       string
	Distance   int
	Score      int
	RoundBonus int

	hand    *Hand
	terrain *TerrainStack
	battle  *BattleStack
	defense *DefenseStack
	travel  *TravelStack

	op *ebiten.DrawImageOptions

	input InputHandler
}

func NewPlayer(id int, name string, input InputHandler, startX, startY float64) *Player {
	var op ebiten.DrawImageOptions
	op.GeoM.Scale(.10, .10)
	op.GeoM.Translate(startX, startY)
	return &Player{
		ID:      id,
		Name:    name,
		hand:    NewHand(),
		terrain: NewTerrainStack(op),
		battle:  NewBattleStack(op),
		defense: NewDefenseStack(op),
		travel:  NewTravelStack(op),
		op:      &op,
		input:   input,
	}
}

func (p *Player) Read() Input {
	return p.input.Read()
}

// Play plays the selected card from player hand into field
func (p *Player) Play() (Card, error) {
	card := p.hand.Remove()
	if card == nil {
		return Card{}, errors.New("no card in hand")
	}
	return *card, nil
}

// Receive receives a card from self (hand) or another player and put on the player field
func (p *Player) Receive(from *Player, card Card) error {
	// validate constraints
	// constraint type: target
	if self := p == from; from != nil && card.Constraints.Target != "" && (self && card.Constraints.Target != "self" || !self && card.Constraints.Target == "self") {
		return fmt.Errorf("can't play card %s on %s field: invalid target", card.Name, p.Name)
	}

	// constraint type: status
	if len(card.Constraints.Status) > 0 {
		validStatus := false
		for _, s := range card.Constraints.Status {
			if s == p.Status() {
				validStatus = true
				break
			}
		}

		if !validStatus {
			return fmt.Errorf("can't play card %s on status %s", card.Name, p.Status())
		}
	}

	// constraint type: terrain
	if len(card.Constraints.Terrains) > 0 && p.terrain.Terrain() != "" {
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
	for _, i := range p.defense.Immunities() {
		if card.Effects.Status == i || (i == StatusLost && card.Type == TypeYellow) {
			return fmt.Errorf("can't play card %s, player has %s immunity", card.Name, i)
		}
	}

	// validate travel rules
	if p.Distance+card.Effects.Distance > config.RoundWin {
		return fmt.Errorf("can't play card because of distance overflow: %d", p.Distance+card.Effects.Distance)
	}

	// add card to proper stack
	switch card.Type {
	case TypeBlue:
		p.defense.Put(card)
	case TypeGreen, TypeRed:
		p.battle.Put(card)
	case TypeWhite:
		p.travel.Put(card)
	case TypeYellow:
		p.terrain.Put(card)
	}

	// calculate effects
	p.Distance += card.Effects.Distance
	p.RoundBonus += card.Effects.Bonus
	if card.counter {
		p.RoundBonus += card.Effects.CounterBonus
	}

	// new immunity
	if p.battle.Status() == card.Effects.Immunity {
		p.battle.Remove()
	}

	// Apply RA terrain effect
	if card.Effects.Immunity == StatusLost {
		p.terrain.Clear()
	}

	return nil
}

func (p *Player) Status() Status {
	status := p.battle.Status()
	if status == StatusUnknown {
		status = StatusLost
	}

	// check if has lost immunity
	for _, i := range p.defense.Immunities() {
		if i == StatusLost {
			switch status {
			case StatusEscaping, StatusHealing, StatusWorking, StatusLost:
				status = StatusOriented
			}
			break
		}
	}

	return status
}

// Defense returns the defense cards in the player field
func (p *Player) Defense() []Card {
	return p.defense.All()
}

// Draw draws a card from the deck and add to the players hand
func (p *Player) DrawCard(d *Deck) error {
	card, err := d.DrawCard()
	if err != nil {
		return err
	}
	p.hand.Put(card)
	return nil
}

// Terrain returns the top card of the terrain stack
func (p *Player) Terrain() Card {
	return p.terrain.Top()
}

// Battle returns the top card of the battle stack
func (p *Player) Battle() Card {
	return p.battle.Top()
}

func (p *Player) Draw(target *ebiten.Image) {
	p.terrain.Draw(target)
	p.battle.Draw(target)
	p.travel.Draw(target)
	p.defense.Draw(target)
}
