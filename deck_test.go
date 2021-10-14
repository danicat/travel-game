package main

import "testing"

func TestInsert(t *testing.T) {
	cards := []Card{{Key: "O"}, {Key: "FD"}, {Key: "SR"}, {Key: "E"}}
	deck := Deck{}
	for _, c := range cards {
		deck.Insert(c)
	}

	if len(deck.cards) != len(cards) {
		t.Fatal("expected deck to have the same number of cards")
	}
}

func TestDraw(t *testing.T) {
	cards := []Card{{Key: "O"}, {Key: "FD"}, {Key: "SR"}, {Key: "E"}}
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

	if card.Key != cards[len(cards)-1].Key {
		t.Fatalf("expected to be the %s card, but got %s", cards[len(cards)-1].Key, card.Key)
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
	cards := []Card{{Key: "O"}, {Key: "FD"}, {Key: "SR"}, {Key: "E"}}
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
		if cards[i].Key != deck.cards[i].Key {
			same = false
			break
		}
	}

	if same {
		t.Fatal("expected deck to be different after shuffle")
	}
}

func TestNewDeck(t *testing.T) {
	cards, err := LoadCards("cards.json")
	if err != nil {
		t.Fatalf("expected no errors, got %s", err)
	}
	d := NewDeck(cards)
	if len(d.cards) != 112 {
		t.Fatal("expected 112 cards")
	}
}
