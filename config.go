package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Rect struct {
	StartX float64 `json:"startX"`
	StartY float64 `json:"startY"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

type Config struct {
	ScreenWidth  int `json:"screen_width"`
	ScreenHeight int `json:"screen_height"`

	MaxPlayers int `json:"max_players"`
	HandSize   int `json:"hand_size"`

	Layout struct {
		Players  []Rect
		System   Rect
		Hand     Rect
		Terrain  Rect
		Battle   Rect
		Travel   Rect
		Deck     Rect
		Cemitery Rect
		Card     Rect
	}
}

var config Config

func LoadConfig(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("unable to open file: %s", err)
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("unable to read file: %s", err)
	}

	err = json.Unmarshal(b, &config)
	if err != nil {
		return fmt.Errorf("error parsing json: %s", err)
	}

	return nil
}
