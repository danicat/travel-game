package main

import (
	"encoding/json"
	"strings"
)

type Status int

const (
	StatusUnknown Status = iota
	StatusOriented
	StatusWorking
	StatusEscaping
	StatusHealing
	StatusLost
	StatusPenniless
	StatusCaptive
	StatusSick
)

func (ps *Status) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	case "oriented":
		*ps = StatusOriented
	case "working":
		*ps = StatusWorking
	case "escaping":
		*ps = StatusEscaping
	case "healing":
		*ps = StatusHealing
	case "lost":
		*ps = StatusLost
	case "penniless":
		*ps = StatusPenniless
	case "captive":
		*ps = StatusCaptive
	case "sick":
		*ps = StatusSick
	}

	return nil
}

func (ps Status) String() string {
	switch ps {
	// green status
	case StatusOriented:
		return "oriented"
	case StatusWorking:
		return "working"
	case StatusEscaping:
		return "escaping"
	case StatusHealing:
		return "healing"
	// red status
	case StatusLost:
		return "lost"
	case StatusPenniless:
		return "penniless"
	case StatusCaptive:
		return "captive"
	case StatusSick:
		return "sick"
	}
	return "unknown"
}
