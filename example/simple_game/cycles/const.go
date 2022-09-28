package main

// type actionType uint

const (
	NoAction uint64 = iota

	StartCompete = iota + 10
	CompeteResule
)

// type competeType uint8

const (
	Rock uint8 = iota
	Paper
	Scissors
)

// type resultType uint8

const (
	Win uint8 = iota
	Lose
	Dogfall
)
