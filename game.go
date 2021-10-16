package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Game struct {
	state    State
	deck     *Deck
	players  *Players
	handSize int

	graveyard []Card
	cardBack  *ebiten.Image

	input InputHandler

	op ebiten.DrawImageOptions
}

func NewGame(maxPlayers, handSize int) (*Game, error) {
	err := LoadCards("cards.json")
	if err != nil {
		return nil, err
	}

	img, _, err := ebitenutil.NewImageFromFile("assets/cards/back.png")
	if err != nil {
		return nil, err
	}

	return &Game{
		players:  NewPlayers(maxPlayers),
		handSize: handSize,
		cardBack: img,
		state:    GameStart,
		input:    NewKBHandler(),
	}, nil
}

func (g *Game) InitDeck() {
	log.Println("initializing deck")
	deck := NewDeck(Cards)
	deck.Shuffle(time.Now().UnixNano())
	g.deck = deck
}

func (g *Game) Deal() {
	log.Println("dealing cards")
	log.Printf("players: %d", g.players.Len())

	for i := 0; i < g.handSize; i++ {
		for _, p := range g.players.All() {
			err := p.Draw(g.deck)
			if err != nil {
				log.Fatalf("error dealing cards: %s", err)
			}
		}
	}
}

func (g *Game) Update() error {
	switch g.state {
	case GameStart:
		g.InitDeck()
		g.Deal()
		g.state = TurnStart
	case TurnStart:
		g.state = Draw
	case Draw:
		if key := g.input.Read(); key == KeySelfOrGraveyard {
			err := g.players.Current().Draw(g.deck)
			if err != nil {
				g.state = GameOver
			}
			g.state = Play
		}
	case Play:
		switch key := g.input.Read(); key {
		case KeyLeft:
			g.players.Current().hand.Left()

		case KeyRight:
			g.players.Current().hand.Right()

		case KeySelfOrGraveyard:
			card, err := g.players.Current().Play()
			if err != nil {
				log.Printf("error playing card %s: %s", card.ID, err)
			}

			g.Play(g.players.Current(), g.players.Current(), card)
			g.state = TurnOver

		case KeyOpponentOrGraveyard:
			card, err := g.players.Current().Play()
			if err != nil {
				log.Printf("error playing card %s: %s", card.ID, err)
			}

			g.Play(g.players.Current(), g.players.PeekNext(), card)
			g.state = TurnOver

		case KeyGraveyard:
			card, err := g.players.Current().Play()
			if err != nil {
				log.Printf("error playing card %s: %s", card.ID, err)
			}

			g.Play(g.players.Current(), nil, card)

			g.state = TurnOver
		}

	case TurnOver:
		g.players.Next()
		g.state = TurnStart
	case GameOver:
		g.input.Cancel()
		os.Exit(0)
	}

	return nil
}

func (g *Game) Play(from *Player, to *Player, card Card) error {
	if to == nil {
		g.graveyard = append(g.graveyard, card)
		return nil
	}

	err := to.Receive(from, card)
	if err != nil {
		log.Println(err)
		g.graveyard = append(g.graveyard, card)
		return err
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	text.Draw(screen, fmt.Sprintf("%s | Distance: %d Phase: %s", g.players.Current().Name, g.players.Current().Distance, g.state), ttfRoboto, 0, config.ScreenHeight-24, color.White)

	g.op.GeoM.Reset()
	g.op.GeoM.Scale(.10, .14)
	g.op.GeoM.Translate(config.Layout.System.StartX+config.Layout.Deck.StartX, config.Layout.System.StartY+config.Layout.Deck.StartY)
	screen.DrawImage(g.cardBack, &g.op)

	if len(g.graveyard) > 0 {
		g.op.GeoM.Reset()
		g.op.GeoM.Scale(.10, .10)
		g.op.GeoM.Translate(config.Layout.System.StartX+config.Layout.Graveyard.StartX, config.Layout.System.StartY+config.Layout.Graveyard.StartY)
		screen.DrawImage(g.graveyard[len(g.graveyard)-1].image, &g.op)
	}

	for i := range g.players.All() {
		if battle := g.players.All()[i].Battle(); battle.ID != "" {
			g.op.GeoM.Reset()
			g.op.GeoM.Scale(.10, .10)
			g.op.GeoM.Translate(config.Layout.Players[i].StartX+config.Layout.Battle.StartX, config.Layout.Players[i].StartY+config.Layout.Battle.StartY)
			screen.DrawImage(battle.image, &g.op)
		}

		if terrain := g.players.All()[i].Terrain(); terrain.ID != "" {
			g.op.GeoM.Reset()
			g.op.GeoM.Scale(.10, .10)
			g.op.GeoM.Translate(config.Layout.Players[i].StartX+config.Layout.Terrain.StartX, config.Layout.Players[i].StartY+config.Layout.Terrain.StartY)
			screen.DrawImage(terrain.image, &g.op)
		}

		g.op.GeoM.Reset()
		g.op.GeoM.Scale(.10, .10)
		g.op.GeoM.Translate(config.Layout.Players[i].StartX+config.Layout.Travel.StartX, config.Layout.Players[i].StartY+config.Layout.Travel.StartY)

		for _, c := range g.players.All()[i].travel.All() {
			screen.DrawImage(c.image, &g.op)
			g.op.GeoM.Translate(20, 0)

		}

		g.op.GeoM.Reset()
		g.op.GeoM.Scale(.10, .10)
		g.op.GeoM.Translate(config.Layout.Players[i].StartX+config.Layout.Defense.StartX, config.Layout.Players[i].StartY+config.Layout.Defense.StartY)

		for _, c := range g.players.All()[i].Defense() {
			screen.DrawImage(c.image, &g.op)
			g.op.GeoM.Translate(config.Layout.Card.Width, 0)
		}
	}

	for i, c := range g.players.Current().Hand() {
		g.op.GeoM.Reset()
		var scale float64
		if g.state == Play && g.players.Current().hand.selected == i {
			scale = .12
		} else {
			scale = .10
		}

		g.op.GeoM.Scale(scale, scale)
		g.op.GeoM.Translate(config.Layout.Hand.StartX+float64(i)*config.Layout.Card.Width, config.Layout.Hand.StartY)
		screen.DrawImage(c.image, &g.op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.ScreenWidth, config.ScreenHeight
}
