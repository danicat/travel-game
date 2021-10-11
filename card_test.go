package main

import (
	"testing"
)

func TestLoadCards(t *testing.T) {
	cards, err := LoadCards("cards.json")
	if err != nil {
		t.Fatalf("expected no errors loading cards, got %s", err)
	}

	if len(cards) == 0 {
		t.Fatal("expected cards to be loaded")
	}

	for _, card := range cards {
		if card.image == nil {
			t.Errorf("expected all cards to have an asset. card %s doesn't have one", card.Name)
		}
	}

	if cards["O"].Type != "green" {
		t.Errorf("%s card should be green, got %s", cards["O"].Name, cards["O"].Type)
	}
}
