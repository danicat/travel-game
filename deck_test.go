package main

import "testing"

func TestInsert(t *testing.T) {
	cards := []string{"O", "FD", "SR", "E"}
	deck := Deck{}
	for _, c := range cards {
		deck.Insert(c)
	}

	if len(deck.cards) != len(cards) {
		t.Fatal("expected deck to have the same number of cards")
	}
}

func TestDraw(t *testing.T) {
	cards := []string{"O", "FD", "SR", "E"}
	deck := Deck{}
	for _, c := range cards {
		deck.Insert(c)
	}

	if len(deck.cards) != len(cards) {
		t.Fatal("expected deck to have the same number of cards")
	}

	card, err := deck.Draw()
	if err != nil {
		t.Fatal("expected no errors drawing card")

	}

	if len(deck.cards) != len(cards)-1 {
		t.Fatal("expected deck to have one less card")
	}

	if card != cards[len(cards)-1] {
		t.Fatalf("expected to be the %s card, but got %s", cards[len(cards)-1], card)
	}
}

func TestDrawNoCards(t *testing.T) {
	deck := Deck{}

	_, err := deck.Draw()
	if err == nil {
		t.Fatal("expected error drawing from empty deck")

	}
}

func TestShuffle(t *testing.T) {
	cards := []string{"O", "FD", "SR", "E"}
	deck := Deck{}
	for _, c := range cards {
		deck.Insert(c)
	}

	if len(deck.cards) != len(cards) {
		t.Fatal("expected deck to have the same number of cards")
	}

	deck.Shuffle(123)

	same := true
	for i := 0; i < len(cards); i++ {
		if cards[i] != deck.cards[i] {
			same = false
			break
		}
	}

	if same {
		t.Fatal("expected deck to be different after shuffle")
	}
}
