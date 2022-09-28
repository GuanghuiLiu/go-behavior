package utils

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandInt(n int) int {
	return rand.Intn(n)
}
func RandUint64() uint64 {
	return rand.Uint64()
}
func RandUint8(n uint8) uint8 {
	return uint8(rand.Intn(int(n)))
}
func RandUint64n(n int) uint64 {
	return uint64(rand.Intn(n))
}
