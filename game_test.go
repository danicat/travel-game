package main

import (
	"fmt"
	"testing"
)

func TestDealPhase(t *testing.T) {
	err := LoadConfig("config.json")
	if err != nil {
		t.Fatal(err)
	}

	game, err := NewGame(config.MaxPlayers, config.HandSize)
	if err != nil {
		t.Fatal(err)
	}

	game.state = GameStart
	game.Update()

	if game.state != TurnStart {
		t.Fatal("game state should transition to TurnStart")
	}

	for i := 0; i < config.MaxPlayers; i++ {
		if len(game.players[i].hand) != config.HandSize {
			t.Fatalf("player %d hand size should be %d, got %d", i, config.HandSize, len(game.players[i].hand))
		}
	}
}

func TestStateTransitions(t *testing.T) {
	err := LoadConfig("config.json")
	if err != nil {
		t.Fatal(err)
	}

	game, err := NewGame(config.MaxPlayers, config.HandSize)
	if err != nil {
		t.Fatal(err)
	}

	tbl := []struct {
		beforePlayer int
		afterPlayer  int
		beforeState  State
		afterState   State
	}{
		{
			0,
			0,
			GameStart,
			TurnStart,
		},
		{
			0,
			0,
			TurnStart,
			Draw,
		},
		{
			0,
			1,
			TurnOver,
			TurnStart,
		},
		{
			1,
			0,
			TurnOver,
			TurnStart,
		},
	}

	for _, testcase := range tbl {
		t.Run(fmt.Sprintf("transition from %s to %s", testcase.beforeState, testcase.afterState), func(t *testing.T) {
			game.currentPlayer = testcase.beforePlayer
			game.state = testcase.beforeState
			game.Update()

			if game.state != testcase.afterState {
				t.Fatalf("expected state %s, got %s", testcase.afterState, game.state)
			}

			if game.currentPlayer != testcase.afterPlayer {
				t.Fatalf("expected player %d, got %d", testcase.afterPlayer, game.currentPlayer)
			}

		})
	}
}
