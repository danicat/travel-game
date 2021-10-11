package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
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

var Cards map[string]Card

func LoadCards(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("unable to open file: %s", err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatalf("unable to read file: %s", err)
	}

	err = json.Unmarshal(b, &Cards)
	if err != nil {
		log.Fatalf("error parsing json: %s", err)
	}

	for k, card := range Cards {
		img, _, err := ebitenutil.NewImageFromFile(card.Asset)
		if err != nil {
			log.Fatalf("error loading card %#v asset: %s", card, err)
		}
		card.image = img
		Cards[k] = card
	}
}
