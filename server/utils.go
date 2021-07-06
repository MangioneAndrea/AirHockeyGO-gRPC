package main

import (
	"fmt"
	"math/rand"
	"time"
)

func generateId() string {
	seededRand := rand.New(
		rand.NewSource(time.Now().UnixNano()))

	return fmt.Sprint(seededRand.Int())
}
