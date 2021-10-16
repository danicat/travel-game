package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	_ "image/png"
)

// Contraints contains the conditions for this card to be allowed in game
// if the conditions aren't met it should be sent to the graveyard
type Constraints struct {
	Target   string // Self, Opponent
	Status   []Status
	Terrains []string
}

type Effects struct {
	Immunity     Status
	Status       Status
	Terrain      string
	Distance     int
	Bonus        int
	CounterBonus int
}

type Card struct {
	ID    string   `json:"id"`
	Name  string   `json:"name"`
	Type  CardType `json:"type"`
	Asset string   `json:"asset"`
	Count int      `json:"count"`

	Effects     `json:"effects"`
	Constraints `json:"constraints"`

	counter bool
	image   *ebiten.Image
}

var Cards []Card

func LoadCards(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("unable to open file: %s", err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("unable to read file: %s", err)
	}

	err = json.Unmarshal(b, &Cards)
	if err != nil {
		return fmt.Errorf("error parsing json: %s", err)
	}

	for k, card := range Cards {
		img, _, err := ebitenutil.NewImageFromFile(card.Asset)
		if err != nil {
			return fmt.Errorf("error loading card %#v asset: %s", card, err)
		}
		card.image = img
		Cards[k] = card
	}

	return nil
}

func (c *Card) String() string {
	return c.ID
}

func FindCardByID(id string) *Card {
	for _, c := range Cards {
		if c.ID == id {
			return &c
		}
	}
	return nil
}
