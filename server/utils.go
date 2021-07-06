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

const DEBUG = true

func printDebug(msg string, props ...interface{}) {
	if DEBUG {
		fmt.Printf(msg, props...)
	}
}
