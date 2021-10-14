package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const HandSize = 6

type Game struct {
	state   State
	deck    *Deck
	cards   map[string]Card
	players []*Player

	cardBack *ebiten.Image

	// Gane State
	currentPlayer int
	cardSelected  int

	op      ebiten.DrawImageOptions
	timeout sync.Once
}

func NewGame() (*Game, error) {
	cards, err := LoadCards("cards.json")
	if err != nil {
		return nil, err
	}

	players := []*Player{
		{},
		{},
	}

	img, _, err := ebitenutil.NewImageFromFile("assets/cards/back.png")
	if err != nil {
		return nil, err
	}

	return &Game{cards: cards, state: GameStart, players: players, cardBack: img}, nil
}

func (g *Game) InitDeck() {
	log.Println("initializing deck")
	deck := NewDeck(g.cards)
	deck.Shuffle(time.Now().UnixNano())
	g.deck = deck
}

func (g *Game) Deal(handSize int) {
	log.Println("dealing cards")
	log.Printf("players: %d", len(g.players))
	for i := 0; i < handSize; i++ {
		for _, p := range g.players {
			err := p.Draw(g.deck)
			if err != nil {
				log.Fatalf("error dealing cards: %s", err)
			}
		}
	}
	log.Printf("player 0 hand: %v", g.players[0].hand)
	log.Printf("player 1 hand: %v", g.players[1].hand)
}

func (g *Game) Update() error {
	g.timeout.Do(func() {
		time.AfterFunc(time.Second*60, func() {
			g.state = GameOver
		})
	})
	switch g.state {
	case GameStart:
		g.InitDeck()
		g.Deal(HandSize)
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

			err = g.Play(card, g.currentPlayer)
			if err != nil {
				log.Printf("error playing card %s: %s", card.Key, err)
			}

			g.state = TurnOver
		}
	case TurnOver:
		g.currentPlayer = g.currentPlayer + 1
		if g.currentPlayer > len(g.players)-1 {
			g.currentPlayer = 0
		}
		g.cardSelected = 0
		g.state = TurnStart
	case GameOver:
		os.Exit(0)
	}
	// time.Sleep(time.Millisecond * 100)
	return nil
}

func (g *Game) Play(card Card, target int) error {
	switch card.Type {
	case "green":
		for _, p := range card.Playable {
			if p == g.players[target].BattleStatus() {
				g.players[target].Receive(card)
			}
		}
	case "yellow":
		g.players[target].Receive(card)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// screen.Fill(color.Transparent)
	// screen.Fill(color.RGBA{R: 233, G: 212, B: 96, A: 0xaf})
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("Player %d Turn", g.currentPlayer))

	// var cx, cy float64
	// for i, k := range g.players[g.currentPlayer].Hand() {
	// 	var scale float64
	// 	g.op.GeoM.Reset()
	// 	if g.state == Play && g.cardSelected == i {
	// 		scale = .12
	// 	} else {
	// 		scale = .10
	// 	}
	// 	g.op.GeoM.Scale(scale, scale)
	// 	cy = 100
	// 	g.op.GeoM.Translate(cx, cy)
	// 	cx += 1000 * scale
	// 	card := k.image
	// 	screen.DrawImage(card, &g.op)
	// }

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
		//g.op.GeoM.Translate(float64(i)*config.Layout.Card.Width+10, 0)
		screen.DrawImage(c.image, &g.op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.ScreenWidth, config.ScreenHeight
}
