package main

import (
	"github.intel.com/kubernetes/ci-cd-broker/broker"
)

func main() {
	// Getting the broker.
	b := broker.GetBroker()
	// Start launching up all kafka consumers.
	b.RunConsumers()
}
