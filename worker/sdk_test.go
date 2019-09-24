package worker

import (
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewV8Worker_CallMethod(t *testing.T) {
	worker := NewV8Worker()
	outputArgs, outputErr, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), `
function hello() {
	return 1
}
`, "hello", nil)
	require.NoError(t, err)
	require.NoError(t, outputErr)
	require.NotNil(t, outputArgs)

	uint32Value := outputArgs.ArgumentsIterator().NextArguments().Uint32Value()
	require.EqualValues(t, 1, uint32Value)
}

func TestNewV8Worker_CallMethodWithArguments(t *testing.T) {
	t.Skip("not implemented")
	worker := NewV8Worker()
	outputArgs, outputErr, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), `
function hello(a, b) {
	return 1 + a + b
}
`, "hello", ArgsToArgumentArray(2, 3))
	require.NoError(t, err)
	require.NoError(t, outputErr)
	require.NotNil(t, outputArgs)

	uint32Value := outputArgs.ArgumentsIterator().NextArguments().Uint32Value()
	require.EqualValues(t, 6, uint32Value)
}
