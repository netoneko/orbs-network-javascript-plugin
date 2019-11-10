package test

import (
	"github.com/orbs-network/orbs-network-javascript-plugin/test"
	. "github.com/orbs-network/orbs-network-javascript-plugin/worker"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewV8Worker_Address(t *testing.T) {
	sdkHandler := test.AFakeSdkFor([]byte("signer"), []byte("caller"))

	contract := `
import { Address } from "orbs-contract-sdk/v1";

export function getSignerAddress(value) {
	return Address.getSignerAddress()
}

export function getCallerAddress() {
	return Address.getCallerAddress()
}

export function getOwnAddress() {
	return Address.getOwnAddress()
}
`
	worker := newTestWorker(t, sdkHandler, contract)

	// signer address
	bytesValue := worker.callMethodWithoutErrors("getSignerAddress", ArgsToArgumentArray())
	require.EqualValues(t, []byte("signer"), bytesValue.BytesValue())

	// caller address
	bytesValue = worker.callMethodWithoutErrors("getCallerAddress", ArgsToArgumentArray())
	require.EqualValues(t, []byte("caller"), bytesValue.BytesValue())

	// FIXME add own address
}
