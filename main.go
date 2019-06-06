package main

import (
	"gitlab.devtools.intel.com/kubernetes/ci-cd-broker/broker"
)

func main() {
	// Getting the broker.
	b := broker.Get()
	// Launching all kafka consumers.
	b.Run()
}
