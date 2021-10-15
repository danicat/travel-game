package main

import "testing"

func TestInsert(t *testing.T) {
	cards := []Card{{ID: "O"}, {ID: "FD"}, {ID: "SR"}, {ID: "E"}}
	deck := Deck{}
	for _, c := range cards {
		deck.Insert(c)
	}

	if len(deck.cards) != len(cards) {
		t.Fatal("expected deck to have the same number of cards")
	}
}

func TestDraw(t *testing.T) {
	cards := []Card{{ID: "O"}, {ID: "FD"}, {ID: "SR"}, {ID: "E"}}
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

	if card.ID != cards[len(cards)-1].ID {
		t.Fatalf("expected to be the %s card, but got %s", cards[len(cards)-1].ID, card.ID)
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
	cards := []Card{{ID: "O"}, {ID: "FD"}, {ID: "SR"}, {ID: "E"}}
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
		if cards[i].ID != deck.cards[i].ID {
			same = false
			break
		}
	}

	if same {
		t.Fatal("expected deck to be different after shuffle")
	}
}

func TestNewDeck(t *testing.T) {
	err := LoadCards("cards.json")
	if err != nil {
		t.Fatalf("expected no errors, got %s", err)
	}
	d := NewDeck(Cards)
	if len(d.cards) != 112 {
		t.Fatal("expected 112 cards")
	}
}
