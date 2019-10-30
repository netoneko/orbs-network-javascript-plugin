package worker

import (
	"github.com/orbs-network/orbs-network-javascript-plugin/test"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestBenchmarkToken(t *testing.T) {
	owner := []byte("owner")
	receiver := []byte("receiver")
	totalSupply := uint32(10000000)
	amount := uint32(1982)

	sdkHandler := test.AFakeSdkFor(owner, owner)

	contract, _ := ioutil.ReadFile("../contracts/benchmark/token-in-progress.js")

	// _init
	worker := NewV8Worker(sdkHandler)
	outputArgs, outputErr, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), string(contract),
		"_init", ArgsToArgumentArray())
	require.NoError(t, err)
	require.NoError(t, outputErr)

	// totalSupply
	outputArgs, outputErr, err = worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), string(contract),
		"totalSupply", ArgsToArgumentArray())
	require.NoError(t, err)
	require.NoError(t, outputErr)

	uin32Value := outputArgs.ArgumentsIterator().NextArguments().Uint32Value()
	require.EqualValues(t, totalSupply, uin32Value)

	// transfer
	outputArgs, outputErr, err = worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), string(contract),
		"transfer", ArgsToArgumentArray(amount, receiver))
	require.NoError(t, err)
	require.NoError(t, outputErr)

	// receiver balance
	outputArgs, outputErr, err = worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), string(contract),
		"balanceOf", ArgsToArgumentArray(receiver))
	require.NoError(t, err)
	require.NoError(t, outputErr)

	receiverBalance := outputArgs.ArgumentsIterator().NextArguments().Uint32Value()
	require.EqualValues(t, amount, receiverBalance)

	// owner balance
	outputArgs, outputErr, err = worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), string(contract),
		"balanceOf", ArgsToArgumentArray(owner))
	require.NoError(t, err)
	require.NoError(t, outputErr)

	ownerBalance := outputArgs.ArgumentsIterator().NextArguments().Uint32Value()
	require.EqualValues(t, totalSupply-amount, ownerBalance)
}
