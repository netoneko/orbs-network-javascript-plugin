package test

import (
	"github.com/orbs-network/orbs-contract-sdk/go/context"
	. "github.com/orbs-network/orbs-network-javascript-plugin/worker"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/stretchr/testify/require"
	"testing"
)

type testWorkerWrapper struct {
	worker   Worker
	contract string
	t        *testing.T
}

func (w *testWorkerWrapper) callMethodWithoutErrors(methodName string, args *protocol.ArgumentArray) *protocol.Argument {
	outputArgs, outputErr := w.callMethodWithErrors(methodName, args)
	require.NoError(w.t, outputErr)

	return outputArgs
}

func (w *testWorkerWrapper) callMethodWithErrors(methodName string, args *protocol.ArgumentArray) (*protocol.Argument, error) {
	outputArgs, outputErr, err := w.worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), w.contract,
		primitives.MethodName(methodName), args)
	require.NoError(w.t, err)

	return outputArgs.ArgumentsIterator().NextArguments(), outputErr
}

func newTestWorker(t *testing.T, handler context.SdkHandler, contract string) *testWorkerWrapper {
	return &testWorkerWrapper{
		worker:   NewV8Worker(handler),
		contract: contract,
		t:        t,
	}
}
