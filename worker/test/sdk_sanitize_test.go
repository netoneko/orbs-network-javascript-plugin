package test

import (
	"github.com/orbs-network/orbs-network-javascript-plugin/test"
	. "github.com/orbs-network/orbs-network-javascript-plugin/worker"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewV8Worker_Sanitize(t *testing.T) {
	sdkHandler := test.AFakeSdkFor([]byte("signer"), []byte("caller"))

	evalContract := `
export function testEval(value) {
	return eval(value)
}
`
	worker := NewV8Worker(sdkHandler)
	_, _, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), evalContract, "testEval", ArgsToArgumentArray("National Treasure"))
	require.EqualError(t, err, `keyword "eval" is forbidden in smart contract code`)

	v8Worker2Contract := `
export function testV8Worker2(value) {
	V8Worker2.print(value)
}
`
	_, _, err = worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), v8Worker2Contract, "testEval", ArgsToArgumentArray("National Treasure"))
	require.EqualError(t, err, `keyword "V8Worker2" is forbidden in smart contract code`)
}
