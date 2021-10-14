package main

type State int

const (
	Undefined State = iota
	GameStart
	GameOver
	RoundStart
	RoundOver
	TurnStart
	TurnOver
	Draw
	Play
	BeforeTarget
	Target
	Counter
)

func (s State) String() string {
	switch s {
	case Undefined:
		return "Undefined"
	case GameStart:
		return "GameStart"
	case GameOver:
		return "GameOVer"
	case RoundStart:
		return "RoundStart"
	case RoundOver:
		return "RoundOver"
	case TurnStart:
		return "TurnStart"
	case TurnOver:
		return "TurnOver"
	case Draw:
		return "Draw"
	case Play:
		return "Play"
	case Target:
		return "Target"
	default:
		return "Unknown"
	}
}
