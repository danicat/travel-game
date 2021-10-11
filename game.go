package main

import (
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {
	state State
	deck  *Deck
	cards map[string]Card
}

func NewGame() (*Game, error) {
	cards, err := LoadCards("cards.json")
	if err != nil {
		return nil, err
	}

	return &Game{cards: cards, state: GameStart}, nil
}

func (g *Game) InitDeck() {
	log.Println("initializing deck")
	deck := NewDeck(g.cards)
	deck.Shuffle(time.Now().UnixNano())
	g.deck = deck
}

func (g *Game) Update() error {
	switch g.state {
	case GameStart:
		if ebiten.IsKeyPressed(ebiten.KeyEnter) {
			g.InitDeck()
			g.state = RoundStart
		} else if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			g.state = GameOver
		}
	case RoundStart:
		if ebiten.IsKeyPressed(ebiten.KeyEscape) {
			g.state = RoundOver
		}
	case RoundOver:
		g.state = GameOver
	case GameOver:
		os.Exit(0)
	default:
		log.Fatalf("invalid game state: %s", g.state)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case GameStart:
		ebitenutil.DebugPrint(screen, "Game Start")
	case RoundStart:
		ebitenutil.DebugPrint(screen, "Round Start")
	case RoundOver:
		ebitenutil.DebugPrint(screen, "Round Over")
	case GameOver:
		ebitenutil.DebugPrint(screen, "Game Over")
	default:
		log.Fatalf("invalid game state: %s", g.state)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}
