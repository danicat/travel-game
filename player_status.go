package main

type PlayerStatus int

const (
	StatusUnknown PlayerStatus = iota
	StatusOriented
	StatusWorking
	StatusEscaping
	StatusHealing
	StatusLost
	StatusPenniless
	StatusCaptive
	StatusSick
)

func (ps PlayerStatus) String() string {
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
