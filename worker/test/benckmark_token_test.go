package test

import (
	"github.com/orbs-network/orbs-network-javascript-plugin/test"
	. "github.com/orbs-network/orbs-network-javascript-plugin/worker"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"testing"
)

func TestBenchmarkToken(t *testing.T) {
	owner := []byte("owner")
	receiver := []byte("receiver")
	totalSupply := uint64(10000000000)
	amount := uint64(1982)

	sdkHandler := test.AFakeSdkFor(owner, owner)

	contract, _ := ioutil.ReadFile("../../contracts/benchmark/token.js")

	// _init
	worker := newTestWorker(t, sdkHandler, string(contract))
	worker.callMethodWithoutErrors("_init", ArgsToArgumentArray())

	// totalSupply
	totalSupplyValue := worker.callMethodWithoutErrors("totalSupply", ArgsToArgumentArray())
	require.EqualValues(t, totalSupply, totalSupplyValue.Uint64Value())

	// transfer
	//sdkHandler.MockEmitEvent(func(from []byte, to []byte, amount uint64) {}, owner, receiver, amount)
	worker.callMethodWithoutErrors("transfer", ArgsToArgumentArray(amount, receiver))

	// receiver balance
	receiverBalance := worker.callMethodWithoutErrors("balanceOf", ArgsToArgumentArray(receiver))
	require.EqualValues(t, amount, receiverBalance.Uint64Value())

	// owner balance
	ownerBalance := worker.callMethodWithoutErrors("balanceOf", ArgsToArgumentArray(owner))
	require.EqualValues(t, totalSupply-amount, ownerBalance.Uint64Value())

	//sdkHandler.VerifyMocks()
}
