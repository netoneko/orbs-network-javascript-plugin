package worker

import (
	"github.com/orbs-network/orbs-network-javascript-plugin/test"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestNewV8Worker_CallMethod(t *testing.T) {
	sdkHandler := test.AFakeSdk()
	worker := NewV8Worker(sdkHandler)
	outputArgs, outputErr, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), `
function hello() {
	return 1
}
`, "hello", ArgsToArgumentArray())
	require.NoError(t, err)
	require.NoError(t, outputErr)
	require.NotNil(t, outputArgs)

	uint32Value := outputArgs.ArgumentsIterator().NextArguments().Uint32Value()
	require.EqualValues(t, 1, uint32Value)
}

func TestNewV8Worker_CallMethodWithArguments(t *testing.T) {
	sdkHandler := test.AFakeSdk()
	worker := NewV8Worker(sdkHandler)
	outputArgs, outputErr, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), `
function hello(a, b) {
	return 1 + a + b
}
`, "hello", ArgsToArgumentArray(uint32(2), uint32(3)))
	require.NoError(t, err)
	require.NoError(t, outputErr)
	require.NotNil(t, outputArgs)

	uint32Value := outputArgs.ArgumentsIterator().NextArguments().Uint32Value()
	require.EqualValues(t, 6, uint32Value)
}

func TestNewV8Worker_CallMethodWithStringArgument(t *testing.T) {
	sdkHandler := test.AFakeSdk()
	worker := NewV8Worker(sdkHandler)
	outputArgs, outputErr, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), `
function hello(a) {
	return "hello, " + a
}
`, "hello", ArgsToArgumentArray("world"))
	require.NoError(t, err)
	require.NoError(t, outputErr)
	require.NotNil(t, outputArgs)

	stringValue := outputArgs.ArgumentsIterator().NextArguments().StringValue()
	require.EqualValues(t, "hello, world", stringValue)
}
