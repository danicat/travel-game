package main

import (
	"fmt"
	"log"
)

type Players struct {
	players       []*Player
	currentPlayer int
}

func NewPlayers(input InputHandler, numPlayers int) *Players {
	if numPlayers == 0 {
		log.Fatal("cannot start game with zero players")
	}

	var players []*Player
	for i := 0; i < numPlayers; i++ {
		name := fmt.Sprintf("Player %d", i+1)
		var ih InputHandler
		if i == 0 {
			ih = input
		} else {
			ih = NewAIHandler()
		}
		p := NewPlayer(i, name, ih, config.Layout.Players[i].StartX, config.Layout.Players[i].StartY)
		players = append(players, p)
	}
	return &Players{players: players}
}

func (p *Players) Current() *Player {
	return p.players[p.currentPlayer]
}

// Next returns the next player in the turn cycle and advance the current player index
func (p *Players) Next() *Player {
	if p.currentPlayer == len(p.players)-1 {
		p.currentPlayer = 0
	} else {
		p.currentPlayer++
	}
	return p.players[p.currentPlayer]
}

// PeekNext returns the next player in the turn cycle but doesn't advance the current player index
func (p *Players) PeekNext() *Player {
	if p.currentPlayer == len(p.players)-1 {
		return p.players[0]
	}
	return p.players[p.currentPlayer+1]
}

func (p *Players) Len() int {
	return len(p.players)
}

func (p *Players) All() []*Player {
	return p.players
}
