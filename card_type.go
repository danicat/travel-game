package main

import (
	"encoding/json"
	"strings"
)

type CardType int

const (
	TypeUnknown CardType = iota
	TypeBlue
	TypeGreen
	TypeYellow
	TypeRed
	TypeWhite
)

func (ct *CardType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	case "blue":
		*ct = TypeBlue
	case "green":
		*ct = TypeGreen
	case "red":
		*ct = TypeRed
	case "white":
		*ct = TypeWhite
	case "yellow":
		*ct = TypeYellow
	default:
		*ct = TypeUnknown
	}

	return nil
}
