package test

import (
	"github.com/orbs-network/orbs-network-javascript-plugin/test"
	. "github.com/orbs-network/orbs-network-javascript-plugin/worker"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewV8Worker_ManipulateStateWithDifferentDataTypes(t *testing.T) {
	sdkHandler := test.AFakeSdkFor([]byte("signer"), []byte("caller"))

	contract := `
import { State } from "orbs-contract-sdk/v1";
const BYTES_KEY = new Uint8Array([1])
const UINT32_KEY = new Uint8Array([2])
const STRING_KEY = new Uint8Array([3])

function writeBytes(value) {
	State.writeBytes(BYTES_KEY, value)
}

function readBytes() {
	return State.readBytes(BYTES_KEY)
}

function writeUint32(value) {
	State.writeUint32(UINT32_KEY, value)
}

function readUint32() {
	return State.readUint32(UINT32_KEY)
}

function writeString(value) {
	State.writeString(STRING_KEY, value)
}

function readString() {
	return State.readString(STRING_KEY)
}
`
	worker := newTestWorker(t, sdkHandler, contract)

	// bytes
	worker.callMethodWithoutErrors("writeBytes", ArgsToArgumentArray([]byte("Diamond Dogs")))
	bytesValue := worker.callMethodWithoutErrors("readBytes", ArgsToArgumentArray())
	require.EqualValues(t, []byte("Diamond Dogs"), bytesValue.BytesValue())

	// uint32
	worker.callMethodWithoutErrors("writeUint32", ArgsToArgumentArray(uint32(123456)))
	uint32Value := worker.callMethodWithoutErrors("readUint32", ArgsToArgumentArray())
	require.EqualValues(t, uint32(123456), uint32Value.Uint32Value())

	// string
	worker.callMethodWithoutErrors("writeString", ArgsToArgumentArray("Diamond Dogs"))
	stringValue := worker.callMethodWithoutErrors("readString", ArgsToArgumentArray())
	require.EqualValues(t, "Diamond Dogs", stringValue.StringValue())

}
