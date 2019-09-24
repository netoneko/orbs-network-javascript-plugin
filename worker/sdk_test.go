package worker

import (
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewV8Worker(t *testing.T) {
	worker := NewV8Worker()
	outputArgs, outputErr, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), `
function hello() {
	return 1
}
`, "hello", nil)
	require.NoError(t, err)
	require.NoError(t, outputErr)
	require.NotNil(t, outputArgs)

	bytesValue := outputArgs.ArgumentsIterator().NextArguments().BytesValue()
	require.EqualValues(t, []byte{1, 0, 0, 0, 0, 0, 0, 0}, bytesValue)
}
