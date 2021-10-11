package main

import (
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/looplab/fsm"
)

type Game struct {
	fsm   *fsm.FSM
	deck  *Deck
	cards map[string]Card
}

func NewGame() (*Game, error) {
	cards, err := LoadCards("cards.json")
	if err != nil {
		return nil, err
	}

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
	return &Game{
		fsm:   fsm,
		cards: cards,
	}, nil
}

func (g *Game) InitDeck() {
	deck := NewDeck(g.cards)
	deck.Shuffle(time.Now().UnixNano())
	g.deck = deck
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
