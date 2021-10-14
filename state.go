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
	default:
		return "Unknown"
	}
}
