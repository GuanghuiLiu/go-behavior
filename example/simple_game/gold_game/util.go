package main

import (
	"math/rand"
	"time"
)

func randUint8(n uint8) uint8 {
	return uint8(randInt(int(n)))
}
func randInt(n int) int {
	return rand.New(rand.NewSource(time.Now().UnixNano())).Intn(n)
}

func isAtString(s string, sa []string) bool {
	for _, v := range sa {
		if s == v {
			return true
		}
	}
	return false
}
func isAtInt(s int, sa []int) bool {
	for _, v := range sa {
		if s == v {
			return true
		}
	}
	return false
}
