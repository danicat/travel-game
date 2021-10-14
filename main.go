package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game, err := NewGame()
	if err != nil {
		log.Fatal(err)
	}

	config, err = LoadConfig("config.json")
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("config: %#v", config)

	ebiten.SetWindowSize(config.ScreenWidth, config.ScreenHeight)
	ebiten.SetWindowTitle("Around the World")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
