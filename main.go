package main

import (
	"github.com/juliocri/ci-cd-broker/broker"
)

func main() {
	// Getting the broker.
	b := broker.Get()
	// Launching all kafka consumers.
	b.Run()
}
