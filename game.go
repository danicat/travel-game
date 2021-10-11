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
	var game Game
	cards, err := LoadCards("cards.json")
	if err != nil {
		return nil, err
	}
	game.cards = cards

	fsm := fsm.NewFSM(
		"GameStart",
		fsm.Events{
			{Name: "GameEnded", Src: []string{"GameStart", "RoundStart", "RoundOver"}, Dst: "GameOver"},
			{Name: "RoundStarted", Src: []string{"GameStart"}, Dst: "RoundStart"},
			{Name: "RoundEnded", Src: []string{"RoundStart"}, Dst: "RoundOver"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) {
				log.Printf("event e: %v", e)
			},
			"enter_RoundStart": func(e *fsm.Event) {
				game.InitDeck()
			},
		},
	)
	game.fsm = fsm

	return &game, nil
}

func (g *Game) InitDeck() {
	log.Println("initializing deck")
	deck := NewDeck(g.cards)
	deck.Shuffle(time.Now().UnixNano())
	g.deck = deck
}

func (g *Game) Update() error {
	switch g.fsm.Current() {
	case "GameStart":
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.fsm.Event("RoundStarted")
		} else if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			g.fsm.Event("GameEnded")
		}
	case "RoundStart":
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			g.fsm.Event("RoundEnded")
		}
	case "RoundOver":
		g.fsm.Event("GameEnded")
	case "GameOver":
		os.Exit(0)
	default:
		log.Fatalf("invalid game state: %s", g.fsm.Current())
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.fsm.Current() {
	case "GameStart":
		ebitenutil.DebugPrint(screen, "Game Start")
	case "RoundStart":
		ebitenutil.DebugPrint(screen, "Round Start")
	case "RoundOver":
		ebitenutil.DebugPrint(screen, "Round Over")
	case "GameOver":
		ebitenutil.DebugPrint(screen, "Game Over")
	default:
		log.Fatalf("invalid game state: %s", g.fsm.Current())
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
