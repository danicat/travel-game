package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const HandSize = 6

type Game struct {
	state    State
	deck     *Deck
	cards    map[string]Card
	players  []*Player
	handSize int

	cemitery []Card
	cardBack *ebiten.Image

	// Gane State
	currentPlayer int
	target        int
	cardSelected  int
	currentCard   Card

	op ebiten.DrawImageOptions
	// timeout sync.Once
}

func NewGame(maxPlayers, handSize int) (*Game, error) {
	var game Game
	cards, err := LoadCards("cards.json")
	if err != nil {
		return nil, err
	}
	game.cards = cards

	for i := 0; i < maxPlayers; i++ {
		game.players = append(game.players, &Player{})
	}

	game.handSize = handSize

	img, _, err := ebitenutil.NewImageFromFile("assets/cards/back.png")
	if err != nil {
		return nil, err
	}
	game.cardBack = img

	return &game, nil
}

func (g *Game) InitDeck() {
	log.Println("initializing deck")
	deck := NewDeck(g.cards)
	deck.Shuffle(time.Now().UnixNano())
	g.deck = deck
}

func (g *Game) Deal() {
	log.Println("dealing cards")
	log.Printf("players: %d", len(g.players))
	for i := 0; i < g.handSize; i++ {
		for _, p := range g.players {
			err := p.Draw(g.deck)
			if err != nil {
				log.Fatalf("error dealing cards: %s", err)
			}
		}
	}
}

func (g *Game) Update() error {
	// g.timeout.Do(func() {
	// 	time.AfterFunc(time.Second*60, func() {
	// 		g.state = GameOver
	// 	})
	// })
	switch g.state {
	case GameStart:
		g.InitDeck()
		g.Deal()
		g.state = TurnStart
	case TurnStart:
		g.state = Draw
	case Draw:
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			err := g.players[g.currentPlayer].Draw(g.deck)
			if err != nil {
				g.state = GameOver
			}
			g.state = Play
		}
	case Play:
		if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
			if g.cardSelected <= 0 {
				g.cardSelected = len(g.players[g.currentPlayer].hand) - 1
			} else {
				g.cardSelected--
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
			if g.cardSelected < len(g.players[g.currentPlayer].hand)-1 {
				g.cardSelected++
			} else {
				g.cardSelected = 0
			}
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			card, err := g.players[g.currentPlayer].Play(g.cardSelected)
			if err != nil {
				log.Printf("error playing card %s: %s", card.Key, err)
			}
			g.currentCard = card

			g.state = BeforeTargetSelection
		}
	case BeforeTargetSelection:
		switch g.currentCard.Type {
		case "red":
			g.target = (g.currentPlayer + 1) % 2
		default:
			g.target = g.currentPlayer
		}
		g.state = TargetSelection

	case TargetSelection:
		if inpututil.IsKeyJustPressed(ebiten.KeyUp) || inpututil.IsKeyJustPressed(ebiten.KeyDown) {
			switch g.currentCard.Type {
			case "yellow":
				g.target = (g.target + 1) % 2
			case "red":
				// oponent x graveyard
				if g.target == -1 {
					g.target = (g.currentPlayer + 1) % 2
				} else {
					g.target = -1
				}

			default:
				// self x graveyard
				if g.target == g.currentPlayer {
					g.target = -1
				} else {
					g.target = g.currentPlayer
				}
			}
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
			err := g.Play(g.currentCard, g.target)
			if err != nil {
				log.Printf("error playing card %s: %s", g.currentCard.Key, err)
			}
			g.state = TurnOver
		}

	case TurnOver:
		g.currentPlayer = (g.currentPlayer + 1) % 2
		g.cardSelected = 0
		g.state = TurnStart
	case GameOver:
		os.Exit(0)
	}
	// time.Sleep(time.Millisecond * 100)
	return nil
}

func (g *Game) Play(card Card, target int) error {
	if target >= 0 {
		targetPlayer := g.players[target]
		switch card.Type {
		case "white":
			if targetPlayer.BattleStatus() == "orientation" {
				if len(targetPlayer.terrain) == 0 {
					return targetPlayer.Receive(card)
				}
				for _, a := range card.Allowed {
					if a == targetPlayer.Terrain().Terrain {
						return targetPlayer.Receive(card)
					}
				}
			}

		case "green":
			for _, p := range card.Playable {
				if p == targetPlayer.BattleStatus() {
					return targetPlayer.Receive(card)
				}
			}
		case "red":
			if targetPlayer.BattleStatus() == "orientation" {
				return targetPlayer.Receive(card)
			}
		case "yellow":
			return targetPlayer.Receive(card)
		}
	}
	g.cemitery = append(g.cemitery, card)
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	text.Draw(screen, fmt.Sprintf("Player %d Phase: %s", g.currentPlayer+1, g.state), ttfRoboto, config.ScreenWidth-300, config.ScreenHeight-30, color.White)

	g.op.GeoM.Reset()
	g.op.GeoM.Scale(.10, .14)
	g.op.GeoM.Translate(config.Layout.System.StartX+config.Layout.Deck.StartX, config.Layout.System.StartY+config.Layout.Deck.StartY)
	screen.DrawImage(g.cardBack, &g.op)

	if len(g.cemitery) > 0 {
		g.op.GeoM.Reset()
		g.op.GeoM.Scale(.10, .10)
		g.op.GeoM.Translate(config.Layout.System.StartX+config.Layout.Cemitery.StartX, config.Layout.System.StartY+config.Layout.Cemitery.StartY)
		screen.DrawImage(g.cemitery[len(g.cemitery)-1].image, &g.op)
	}

	for i := range g.players {
		if battle := g.players[i].Battle(); battle.Key != "" {
			g.op.GeoM.Reset()
			g.op.GeoM.Scale(.10, .10)
			g.op.GeoM.Translate(config.Layout.Players[i].StartX+config.Layout.Battle.StartX, config.Layout.Players[i].StartY+config.Layout.Battle.StartY)
			screen.DrawImage(battle.image, &g.op)
		}

		if terrain := g.players[i].Terrain(); terrain.Key != "" {
			g.op.GeoM.Reset()
			g.op.GeoM.Scale(.10, .10)
			g.op.GeoM.Translate(config.Layout.Players[i].StartX+config.Layout.Terrain.StartX, config.Layout.Players[i].StartY+config.Layout.Terrain.StartY)
			screen.DrawImage(terrain.image, &g.op)
		}

		g.op.GeoM.Reset()
		g.op.GeoM.Scale(.10, .10)
		g.op.GeoM.Translate(config.Layout.Players[i].StartX+config.Layout.Travel.StartX, config.Layout.Players[i].StartY+config.Layout.Travel.StartY)

		for _, c := range g.players[i].travel {
			screen.DrawImage(c.image, &g.op)
			g.op.GeoM.Translate(20, 0)

		}
	}

	for i, c := range g.players[g.currentPlayer].Hand() {
		g.op.GeoM.Reset()
		var scale float64
		if g.state == Play && g.cardSelected == i {
			scale = .12
		} else {
			scale = .10
		}

		g.op.GeoM.Scale(scale, scale)
		g.op.GeoM.Translate(config.Layout.Hand.StartX+float64(i)*config.Layout.Card.Width, config.Layout.Hand.StartY)
		screen.DrawImage(c.image, &g.op)
	}

	switch g.state {
	case TargetSelection:
		var choice string
		switch g.target {
		case -1:
			choice = "graveyard"
		case g.currentPlayer:
			choice = "self"
		default:
			choice = "oponent"
		}
		text.Draw(screen, choice, ttfRoboto, config.ScreenWidth/2, config.ScreenHeight/2, color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.ScreenWidth, config.ScreenHeight
}
