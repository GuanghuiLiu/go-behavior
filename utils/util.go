package utils

import (
	"time"
)

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func Min(a, b int) int {
	if a > b {
		return b
	}
	return a
}

func MaxTime(a, b time.Duration) time.Duration {
	if a > b {
		return a
	}
	return b
}
func MinTime(a, b time.Duration) time.Duration {
	if a > b {
		return b
	}
	return a
}
