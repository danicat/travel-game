package main

import (
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/looplab/fsm"
)

type Game struct {
	fsm *fsm.FSM
}

func NewGame() *Game {
	fsm := fsm.NewFSM(
		"GameStart",
		fsm.Events{
			{Name: "GameEnded", Src: []string{"GameStart"}, Dst: "GameOver"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) {
				log.Printf("event e: %v", e)
			},
		},
	)
	return &Game{fsm: fsm}
}

func (g *Game) Update() error {
	switch g.fsm.Current() {
	case "GameStart":
		g.fsm.Event("GameEnded")
	case "GameOver":
		os.Exit(0)
	default:
		log.Fatalf("invalid game state: %s", g.fsm.Current())
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Hello, World!")
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
