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

	game, err := NewGame(nil, config.MaxPlayers, config.HandSize)
	if err != nil {
		t.Fatal(err)
	}

	game.state = RoundStart
	game.Update()

	if game.state != TurnStart {
		t.Fatal("game state should transition to TurnStart")
	}

	for i := 0; i < config.MaxPlayers; i++ {
		if len(game.players.All()[i].hand.All()) != config.HandSize {
			t.Fatalf("player %d hand size should be %d, got %d", i, config.HandSize, len(game.players.All()[i].hand.All()))
		}
	}
}

func TestStateInputlessTransitions(t *testing.T) {
	err := LoadConfig("config.json")
	if err != nil {
		t.Fatal(err)
	}

	game, err := NewGame(nil, config.MaxPlayers, config.HandSize)
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
			RoundStart,
		},
		{
			0,
			0,
			RoundStart,
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
		{
			0,
			0,
			AfterRoundOver,
			GameOver,
		},
	}

	for _, testcase := range tbl {
		t.Run(fmt.Sprintf("transition from %s to %s", testcase.beforeState, testcase.afterState), func(t *testing.T) {
			game.players.Current().ID = testcase.beforePlayer
			game.state = testcase.beforeState
			game.Update()

			if game.state != testcase.afterState {
				t.Fatalf("expected state %s, got %s", testcase.afterState, game.state)
			}

			if game.players.Current().ID != testcase.afterPlayer {
				t.Fatalf("expected player %d, got %d", testcase.afterPlayer, game.players.Current().ID)
			}
		})
	}
}

func TestStateTransitions(t *testing.T) {
	err := LoadConfig("config.json")
	if err != nil {
		t.Fatal(err)
	}

	tbl := []struct {
		name        string
		beforeState State
		inputs      []Input
		afterState  State
	}{
		{
			"should draw a card",
			Draw,
			[]Input{KeyDefaultOrGraveyard},
			Play,
		},
		{
			"should play a card",
			Play,
			[]Input{KeyDefaultOrGraveyard},
			TurnOver,
		},
	}

	for _, testcase := range tbl {
		t.Run(testcase.name, func(t *testing.T) {
			input := NewMockHandler()

			game, err := NewGame(input, config.MaxPlayers, config.HandSize)
			if err != nil {
				t.Fatal(err)
			}

			input.AppendKeys(testcase.inputs)
			game.state = testcase.beforeState
			game.InitDeck()
			game.Deal()
			game.Update()

			if game.state != testcase.afterState {
				t.Fatalf("expected state %s, got %s", testcase.afterState, game.state)
			}
		})
	}
}
