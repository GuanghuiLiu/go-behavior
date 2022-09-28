package main

import (
	"math/rand"
	"time"
)

func randCompete() uint8 {
	return Rock + randUint8(Scissors-Rock+1)
}

func randUint8(n uint8) uint8 {
	return uint8(randInt8(int(n)))
}
func randInt8(n int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(n)
}
