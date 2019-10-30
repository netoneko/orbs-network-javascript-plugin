package test

import (
	"github.com/orbs-network/orbs-contract-sdk/go/context"
	"github.com/orbs-network/orbs-network-javascript-plugin/test"
	. "github.com/orbs-network/orbs-network-javascript-plugin/worker"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewV8Worker_CallSDKHandlerMethod(t *testing.T) {
	sdkHandler := test.AFakeSdkFor([]byte("signer"), []byte("caller"))

	expectedAddr := sdkHandler.SdkAddressGetSignerAddress([]byte("test"), context.PERMISSION_SCOPE_SERVICE)
	require.EqualValues(t, []byte("signer"), expectedAddr)

	contract := `
import { Address } from "orbs-contract-sdk/v1";
export function testSignerAddress(a, b, c) {
	return Address.getSignerAddress()
}
`
	worker := newTestWorker(t, sdkHandler, contract)
	bytesValue := worker.callMethodWithoutErrors("testSignerAddress", ArgsToArgumentArray())
	require.EqualValues(t, []byte("signer"), bytesValue.BytesValue())
}
