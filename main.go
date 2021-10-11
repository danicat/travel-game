package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	// load config

	// load assets
	LoadCards("cards.json")

	// initialize game object
	game := NewGame()

	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Around the World")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
