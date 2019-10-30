package test

import (
	"github.com/orbs-network/orbs-network-javascript-plugin/test"
	. "github.com/orbs-network/orbs-network-javascript-plugin/worker"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewV8Worker_CallMethodWithNoArguments(t *testing.T) {
	sdkHandler := test.AFakeSdk()
	contract := `
function hello() {
	return 1
}`
	worker := newTestWorker(t, sdkHandler, contract)
	uint32Value := worker.callMethodWithoutErrors("hello", ArgsToArgumentArray())
	require.EqualValues(t, 1, uint32Value.Uint32Value())
}

func TestNewV8Worker_CallMethodWithArguments(t *testing.T) {
	sdkHandler := test.AFakeSdk()
	contract := `
function hello(a, b) {
	return 1 + a + b
}
`
	worker := newTestWorker(t, sdkHandler, contract)
	uint32Value := worker.callMethodWithoutErrors("hello", ArgsToArgumentArray(uint32(2), uint32(3)))
	require.EqualValues(t, 6, uint32Value.Uint32Value())
}
