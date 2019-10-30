package worker

import (
	"github.com/orbs-network/orbs-network-javascript-plugin/test"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestBenchmarkToken(t *testing.T) {
	sdkHandler := test.AFakeSdkFor([]byte("signer"), []byte("caller"))

	contract, _ := ioutil.ReadFile("../contracts/benchmark/token.js")

	worker := NewV8Worker(sdkHandler)
	outputArgs, outputErr, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), string(contract),
		"_init", ArgsToArgumentArray())
	require.NoError(t, err)
	require.NoError(t, outputErr)

	outputArgs, outputErr, err = worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), string(contract),
		"totalSupply", ArgsToArgumentArray())
	require.NoError(t, err)
	require.NoError(t, outputErr)

	uin32Value := outputArgs.ArgumentsIterator().NextArguments().Uint32Value()
	require.EqualValues(t, uint32(10000000), uin32Value)
}
