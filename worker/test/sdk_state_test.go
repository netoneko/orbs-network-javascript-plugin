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

export function writeBytes(value) {
	State.writeBytes(BYTES_KEY, value)
}

export function readBytes() {
	return State.readBytes(BYTES_KEY)
}

export function writeUint32(value) {
	State.writeUint32(UINT32_KEY, value)
}

export function readUint32() {
	return State.readUint32(UINT32_KEY)
}

export function writeUint64(value) {
	State.writeUint64(UINT32_KEY, value)
}

export function readUint64() {
	return State.readUint64(UINT32_KEY)
}

export function writeString(value) {
	State.writeString(STRING_KEY, value)
}

export function readString() {
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

	// uint64
	worker.callMethodWithoutErrors("writeUint64", ArgsToArgumentArray(uint64(1234567890123)))
	uint64Value := worker.callMethodWithoutErrors("readUint64", ArgsToArgumentArray())
	require.EqualValues(t, uint64(1234567890123), uint64Value.Uint64Value())

	// string
	worker.callMethodWithoutErrors("writeString", ArgsToArgumentArray("Diamond Dogs"))
	stringValue := worker.callMethodWithoutErrors("readString", ArgsToArgumentArray())
	require.EqualValues(t, "Diamond Dogs", stringValue.StringValue())

}

func TestNewV8Worker_ReadDefaultValuesFromState(t *testing.T) {
	sdkHandler := test.AFakeSdkFor([]byte("signer"), []byte("caller"))

	contract := `
import { State } from "orbs-contract-sdk/v1";
const KEY = new Uint8Array([1])

export function readBytes() {
	return State.readBytes(KEY)
}

export function readUint32() {
	return State.readUint32(KEY)
}

export function readUint64() {
	return State.readUint64(KEY)
}

export function readString() {
	return State.readString(KEY)
}
`
	worker := newTestWorker(t, sdkHandler, contract)

	// bytes
	bytesValue := worker.callMethodWithoutErrors("readBytes", ArgsToArgumentArray())
	require.EqualValues(t, []byte{}, bytesValue.BytesValue())

	// uint32
	uint32Value := worker.callMethodWithoutErrors("readUint32", ArgsToArgumentArray())
	require.EqualValues(t, uint32(0), uint32Value.Uint32Value())

	// uint64
	uint64Value := worker.callMethodWithoutErrors("readUint64", ArgsToArgumentArray())
	require.EqualValues(t, uint64(0), uint64Value.Uint64Value())

	// string
	stringValue := worker.callMethodWithoutErrors("readString", ArgsToArgumentArray())
	require.EqualValues(t, "", stringValue.StringValue())

}

func TestNewV8Worker_ClearState(t *testing.T) {
	sdkHandler := test.AFakeSdkFor([]byte("signer"), []byte("caller"))

	contract := `
import { State } from "orbs-contract-sdk/v1";
const BYTES_KEY = new Uint8Array([1])

export function writeBytes(value) {
	State.writeBytes(BYTES_KEY, value)
}

export function readBytes() {
	return State.readBytes(BYTES_KEY)
}

export function clearBytes() {
	State.clear(BYTES_KEY)
}
`
	worker := newTestWorker(t, sdkHandler, contract)

	worker.callMethodWithoutErrors("writeBytes", ArgsToArgumentArray([]byte{1, 2, 3, 4, 5}))
	bytesValue := worker.callMethodWithoutErrors("readBytes", ArgsToArgumentArray())
	require.EqualValues(t, []byte{1, 2, 3, 4, 5}, bytesValue.BytesValue())

	worker.callMethodWithoutErrors("clearBytes", ArgsToArgumentArray())

	bytesValue = worker.callMethodWithoutErrors("readBytes", ArgsToArgumentArray())
	require.EqualValues(t, []byte{}, bytesValue.BytesValue())

}
