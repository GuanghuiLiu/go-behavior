package main

// type actionType uint

// proto：key
const (
	NoAction uint = iota

	StartGame = iota + 10
	NpcName
	GameInfo
	RoleChoose
	Rolechose
	GameResule
)

const (
	MaxGold       int = 64 * 64
	MaxGameCount  int = 64
	BaseGoldNumbe int = 16
)
