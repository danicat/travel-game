package main

import (
	"testing"
)

func TestLoadCards(t *testing.T) {
	if len(Cards) != 0 {
		t.Fatal("expected no cards loaded before running test")
	}

	LoadCards("cards.json")

	if len(Cards) == 0 {
		t.Fatal("expected cards to be loaded")
	}

	for _, card := range Cards {
		if card.image == nil {
			t.Errorf("expected all cards to have an asset. card %s doesn't have one", card.Name)
		}
	}

	if Cards["O"].Type != "green" {
		t.Errorf("%s card should be green, got %s", Cards["O"].Name, Cards["O"].Type)
	}
}
