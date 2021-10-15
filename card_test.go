package main

import (
	"testing"
)

func TestLoadCards(t *testing.T) {
	err := LoadCards("cards.json")
	if err != nil {
		t.Fatalf("expected no errors loading cards, got %s", err)
	}

	if len(Cards) == 0 {
		t.Fatal("expected cards to be loaded")
	}

	for _, card := range Cards {
		if card.image == nil {
			t.Errorf("expected all cards to have an asset. card %s doesn't have one", card.Name)
		}
	}
}
