package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
)

type Game struct {
	state     State
	deck      *Deck
	players   *Players
	graveyard *Graveyard
	input     InputHandler
	scene     Scene
	handSize  int
	skip      *time.Timer
}

func NewGame(input InputHandler, maxPlayers, handSize int) (*Game, error) {
	err := LoadCards("cards.json")
	if err != nil {
		return nil, err
	}

	return &Game{
		state:     GameStart,
		players:   NewPlayers(input, maxPlayers),
		graveyard: NewGraveyard(),
		input:     input,
		handSize:  handSize,
	}, nil
}

func (g *Game) InitDeck() {
	log.Println("initializing deck")
	deck, err := NewDeck(Cards)
	if err != nil {
		log.Fatal(err)
	}

	deck.Shuffle(time.Now().UnixNano())
	g.deck = deck
}

func (g *Game) Deal() {
	log.Println("dealing cards")
	log.Printf("players: %d", g.players.Len())

	for i := 0; i < g.handSize; i++ {
		for _, p := range g.players.All() {
			err := p.DrawCard(g.deck)
			if err != nil {
				log.Fatalf("error dealing cards: %s", err)
			}
		}
	}
}

func (g *Game) Update() error {
	// Handle state
	switch g.state {
	case GameStart:
		g.state = RoundStart

	case RoundStart:
		g.InitDeck()
		g.Deal()
		g.state = TurnStart

	case TurnStart:
		g.RenderPlayScene()
		g.state = Draw

	case Draw:
		g.players.Current().hand.selected = -1
		switch key := g.players.Current().Read(); key {
		case KeyDefaultOrGraveyard:
			err := g.players.Current().DrawCard(g.deck)
			if err != nil {
				g.state = GameOver
			}
			g.players.Current().hand.selected = len(g.players.Current().hand.All()) - 1
			g.state = Play
		}

	case Play:
		switch key := g.players.Current().Read(); key {
		case KeyLeft:
			g.players.Current().hand.Left()

		case KeyRight:
			g.players.Current().hand.Right()

		case KeyDefaultOrGraveyard, KeyOpponentOrGraveyard:
			card, err := g.players.Current().Play()
			if err != nil {
				// TODO: shouldn't happen? - enable last minute arrival
				g.state = BeforeRoundOver
				break
			}

			target := g.players.Current()
			switch key {
			case KeyDefaultOrGraveyard:
				if card.Type == TypeRed {
					target = g.players.PeekNext()
				}
			case KeyOpponentOrGraveyard:
				target = g.players.PeekNext()
			}

			g.Play(g.players.Current(), target, card)

			// if last card played is blue, player should draw a new card
			// this will force player to play again (in a 2 player game)
			if card.Type == TypeBlue {
				g.players.Next()
			}

			g.state = TurnOver

		case KeyGraveyard:
			card, err := g.players.Current().Play()
			if err != nil {
				log.Printf("error playing card %s: %s", card.ID, err)
			}

			g.Play(g.players.Current(), nil, card)

			g.state = TurnOver

		case KeyQuit:
			g.state = GameOver
		}

	case TurnOver:
		if g.players.Current().Distance == config.RoundWin {
			g.state = BeforeRoundOver
		} else {
			g.players.Next()
			g.state = TurnStart
		}

	case BeforeRoundOver:
		g.skip = time.AfterFunc(time.Second*3, func() {
			log.Println("timer func called")
			g.state = AfterRoundOver
		})
		g.state = RoundOver

	case RoundOver:
		switch key := g.input.Read(); key {
		case KeyDefaultOrGraveyard:
			g.skip.Stop()
			g.state = AfterRoundOver
		}

	case AfterRoundOver:
		log.Println("after round over")
		g.state = GameOver

	case GameOver:
		g.input.Cancel()
		os.Exit(0)
	}

	return nil
}

func (g *Game) RenderPlayScene() {
	g.scene = Scene{}
	g.scene.AddSprite(g.deck)
	g.scene.AddSprite(g.graveyard)
	for _, p := range g.players.All() {
		g.scene.AddSprite(p)
	}
	g.scene.AddSprite(g.players.Current().hand)
}

func (g *Game) Play(from *Player, to *Player, card Card) error {
	if to == nil {
		g.graveyard.Put(card)
		return nil
	}

	err := to.Receive(from, card)
	if err != nil {
		log.Println(err)
		g.graveyard.Put(card)
		return err
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.state {
	case TurnStart, Draw, Play, TurnOver:
		screen.Fill(color.RGBA{R: 00, G: 0x77, B: 0xBE, A: 0xFF})
		text.Draw(screen, fmt.Sprintf("%s | Distance: %d Phase: %s", g.players.Current().Name, g.players.Current().Distance, g.state), ttfRoboto, 0, config.ScreenHeight-24, color.White)

		for _, s := range g.scene.Sprites() {
			s.Draw(screen)
		}
	case BeforeRoundOver, RoundOver, AfterRoundOver:
		screen.Fill(color.RGBA{R: 00, G: 0x77, B: 0xBE, A: 0xFF})
		msg := fmt.Sprintf("Round Over! Player 1: %d | Player 2: %d", g.players.All()[0].Distance, g.players.All()[1].Distance)
		bounds := text.BoundString(ttfRoboto, msg)
		text.Draw(screen, msg, ttfRoboto, config.ScreenWidth/2-bounds.Dx()/2, config.ScreenHeight/2-bounds.Dy()/2, color.White)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.ScreenWidth, config.ScreenHeight
}
