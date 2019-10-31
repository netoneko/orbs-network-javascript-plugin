package test

import (
	"github.com/orbs-network/orbs-network-javascript-plugin/test"
	. "github.com/orbs-network/orbs-network-javascript-plugin/worker"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewV8Worker_MethodNotFound(t *testing.T) {
	sdkHandler := test.AFakeSdkFor([]byte("signer"), []byte("caller"))

	contract := `
import { State } from "orbs-contract-sdk/v1";
const KEY = new Uint8Array([1, 2, 3])

export function write(value) {
	State.writeString(KEY, value)
}
`

	worker := newTestWorker(t, sdkHandler, contract)
	worker.callMethodWithoutErrors("write", ArgsToArgumentArray("Diamond Dogs"))

	outputValue, outputErr := worker.callMethodWithErrors("_read", ArgsToArgumentArray())
	require.EqualError(t, outputErr, "JS contract execution failed")
	require.EqualValues(t, "method '_read' not found in contract", outputValue.StringValue())
}

func TestNewV8Worker_MethodThrowsError(t *testing.T) {
	sdkHandler := test.AFakeSdkFor([]byte("signer"), []byte("caller"))

	contract := `
import { State } from "orbs-contract-sdk/v1";
const KEY = new Uint8Array([1, 2, 3])

export function write(value) {
	State.writeString(KEY, value)
}

export function bang() {
	throw new Error("bang!")
}
`

	worker := newTestWorker(t, sdkHandler, contract)
	worker.callMethodWithoutErrors("write", ArgsToArgumentArray("Diamond Dogs"))

	outputValue, outputErr := worker.callMethodWithErrors("bang", ArgsToArgumentArray())
	require.EqualError(t, outputErr, "JS contract execution failed")
	require.EqualValues(t, "bang!", outputValue.StringValue())
}

func TestNewV8Worker_VerifyDataTypes(t *testing.T) {
	sdkHandler := test.AFakeSdkFor([]byte("signer"), []byte("caller"))

	contract := `
import { Verify } from "orbs-contract-sdk/v1";

export function verifyBytes(value) {
	Verify.bytes(value)
}

export function verifyString(value) {
	Verify.string(value)
}

export function verifyUint32(value) {
	Verify.uint32(value)
}

export function verifyUint64(value) {
	Verify.uint64(value)
}

`

	worker := newTestWorker(t, sdkHandler, contract)

	// bytes
	worker.callMethodWithoutErrors("verifyBytes", ArgsToArgumentArray([]byte("Nicolas Cage")))

	outputValue, outputErr := worker.callMethodWithErrors("verifyBytes", ArgsToArgumentArray("Vampire's Kiss"))
	require.EqualError(t, outputErr, "JS contract execution failed")
	require.EqualValues(t, `Value "Vampire's Kiss" is not a byte array`, outputValue.StringValue())

	// string
	worker.callMethodWithoutErrors("verifyString", ArgsToArgumentArray("Nicolas Cage"))

	outputValue, outputErr = worker.callMethodWithErrors("verifyString", ArgsToArgumentArray([]byte("Vampire's Kiss")))
	require.EqualError(t, outputErr, "JS contract execution failed")
	require.EqualValues(t, `Value "86,97,109,112,105,114,101,39,115,32,75,105,115,115" is not a string`, outputValue.StringValue())

	// uint32
	worker.callMethodWithoutErrors("verifyUint32", ArgsToArgumentArray(uint32(1982)))

	outputValue, outputErr = worker.callMethodWithErrors("verifyUint32", ArgsToArgumentArray(uint64(1997)))
	require.EqualError(t, outputErr, "JS contract execution failed")
	require.EqualValues(t, `Value "1997" is not a uint32`, outputValue.StringValue())

	//// uint64
	worker.callMethodWithoutErrors("verifyUint64", ArgsToArgumentArray(uint64(1997)))

	outputValue, outputErr = worker.callMethodWithErrors("verifyUint64", ArgsToArgumentArray())
	require.EqualError(t, outputErr, "JS contract execution failed")
	require.EqualValues(t, `Value "undefined" is not a uint64`, outputValue.StringValue())
}
