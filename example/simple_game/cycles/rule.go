package main

func rule(self, opponent uint8) (uint8, bool) {
	if self == opponent {
		return Dogfall, true
	}
	switch {
	case self == Rock && opponent == Paper:
		return Lose, true
	case self == Rock && opponent == Scissors:
		return Win, true
	case self == Scissors && opponent == Paper:
		return Win, true
	case self == Scissors && opponent == Rock:
		return Lose, true
	case self == Paper && opponent == Scissors:
		return Lose, true
	case self == Paper && opponent == Rock:
		return Win, true
	}
	return Win, false
}
