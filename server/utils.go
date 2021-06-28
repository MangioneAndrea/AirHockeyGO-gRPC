package main

import (
	"math/rand"
	"time"
)

func generateId() string {
	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))

	return string(seededRand.Int())
}
