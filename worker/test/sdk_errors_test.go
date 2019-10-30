package test

import (
	"github.com/orbs-network/orbs-network-javascript-plugin/test"
	. "github.com/orbs-network/orbs-network-javascript-plugin/worker"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewV8Worker_MethodNotFound(t *testing.T) {
	sdkHandler := test.AFakeSdkFor([]byte("signer"), []byte("caller"))

	contract := `
import { State } from "orbs-contract-sdk/v1";
const KEY = new Uint8Array([1, 2, 3])

function write(value) {
	State.writeString(KEY, value)
}
`

	worker := NewV8Worker(sdkHandler)
	_, outputErr, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), contract,
		"write", ArgsToArgumentArray("Diamond Dogs"))
	require.NoError(t, err)
	require.NoError(t, outputErr)

	outputArgs, outputErr, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), contract,
		"_read", ArgsToArgumentArray())
	require.NoError(t, err)
	require.EqualError(t, outputErr, "JS contract execution failed")
	require.EqualValues(t, "method '_read' not found in contract", outputArgs.ArgumentsIterator().NextArguments().StringValue())
}

func TestNewV8Worker_MethodThrowsError(t *testing.T) {
	sdkHandler := test.AFakeSdkFor([]byte("signer"), []byte("caller"))

	contract := `
import { State } from "orbs-contract-sdk/v1";
const KEY = new Uint8Array([1, 2, 3])

function write(value) {
	State.writeString(KEY, value)
}

function bang() {
	throw new Error("bang!")
}
`

	worker := newTestWorker(t, sdkHandler, contract)
	worker.callMethodWithoutErrors("write", ArgsToArgumentArray("Diamond Dogs"))

	outputValue, outputErr := worker.callMethodWithErrors("bang", ArgsToArgumentArray())
	require.EqualError(t, outputErr, "JS contract execution failed")
	require.EqualValues(t, "bang!", outputValue.StringValue())
}
