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

type Card struct {
	Name        string
	Type        string
	Asset       string
	Effect      string
	Playable    []string
	Counter     string
	PermaEffect string
	Value       int
	Allowed     []string
	Count       int
	image       *ebiten.Image
}

func LoadCards(file string) (map[string]Card, error) {
	var cards map[string]Card

	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("unable to open file: %s", err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("unable to read file: %s", err)
	}

	err = json.Unmarshal(b, &cards)
	if err != nil {
		return nil, fmt.Errorf("error parsing json: %s", err)
	}

	for k, card := range cards {
		img, _, err := ebitenutil.NewImageFromFile(card.Asset)
		if err != nil {
			return nil, fmt.Errorf("error loading card %#v asset: %s", card, err)
		}
		card.image = img
		cards[k] = card
	}

	return cards, nil
}
