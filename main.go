package main

import (
	"github.com/orbs-network/orbs-contract-sdk/go/context"
	"github.com/orbs-network/orbs-network-go/services/processor"
	"github.com/orbs-network/orbs-network-javascript-plugin/worker"
)

func New(handler context.SdkHandler) processor.StatelessProcessor {
	return worker.NewV8Worker(handler)
}