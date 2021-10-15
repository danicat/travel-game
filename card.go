package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"

	_ "image/png"
)

type CardType int

const (
	TypeUnknown CardType = iota
	TypeBlue
	TypeGreen
	TypeYellow
	TypeRed
	TypeWhite
)

func (ct *CardType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	case "blue":
		*ct = TypeBlue
	case "green":
		*ct = TypeGreen
	case "red":
		*ct = TypeRed
	case "white":
		*ct = TypeWhite
	case "yellow":
		*ct = TypeYellow
	default:
		*ct = TypeUnknown
	}

	return nil
}

// Contraints contains the conditions for this card to be allowed in game
// if the conditions aren't met it should be sent to the graveyard
type Constraints struct {
	Terrains []string
	Status   []string
	Target   string // Self, Opponent
}

type Effects struct {
	Immunity     []string
	Status       string
	Terrain      string
	Distance     int
	Bonus        int
	CounterBonus int
}

type Card struct {
	ID    string
	Name  string
	Type  CardType
	Asset string
	Count int

	Effects
	Constraints

	image *ebiten.Image
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
