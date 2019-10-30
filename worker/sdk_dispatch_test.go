package worker

import (
	"github.com/orbs-network/orbs-contract-sdk/go/context"
	"github.com/orbs-network/orbs-network-javascript-plugin/test"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewV8Worker_CallSDKHandlerMethod(t *testing.T) {
	sdkHandler := test.AFakeSdkFor([]byte("signer"), []byte("caller"))

	expectedAddr := sdkHandler.SdkAddressGetSignerAddress([]byte("test"), context.PERMISSION_SCOPE_SERVICE)
	require.EqualValues(t, []byte("signer"), expectedAddr)

	worker := NewV8Worker(sdkHandler)
	outputArgs, outputErr, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), `
import { Address } from "orbs-contract-sdk/v1";
function testSignerAddress(a, b, c) {
	const address = Address.getSignerAddress()
	return address 
}
`, "testSignerAddress", ArgsToArgumentArray(uint32(1), uint32(2), uint32(3)))
	require.NoError(t, err)
	require.NoError(t, outputErr)
	require.NotNil(t, outputArgs)

	bytesValue := outputArgs.ArgumentsIterator().NextArguments().BytesValue()
	require.EqualValues(t, []byte("signer"), bytesValue)
}

func TestNewV8Worker_ManipulateStateWithBytes(t *testing.T) {
	sdkHandler := test.AFakeSdkFor([]byte("signer"), []byte("caller"))

	contract := `
import { State } from "orbs-contract-sdk/v1";
const KEY = new Uint8Array([1, 2, 3, 4, 5])

function write(value) {
	State.writeBytes(KEY, value)
}

function read() {
	return State.readBytes(KEY)
}
`

	worker := NewV8Worker(sdkHandler)
	outputArgs, outputErr, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), contract,
		"write", ArgsToArgumentArray([]byte("Diamond Dogs")))
	require.NoError(t, err)
	require.NoError(t, outputErr)

	outputArgs, outputErr, err = worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), contract,
		"read", ArgsToArgumentArray())
	require.NoError(t, err)
	require.NoError(t, outputErr)

	bytesValue := outputArgs.ArgumentsIterator().NextArguments().BytesValue()
	require.EqualValues(t, []byte("Diamond Dogs"), bytesValue)
}

func TestNewV8Worker_ManipulateStateWithUint32(t *testing.T) {
	sdkHandler := test.AFakeSdkFor([]byte("signer"), []byte("caller"))

	contract := `
import { State } from "orbs-contract-sdk/v1";
const KEY = new Uint8Array([1, 2, 3])

function write(value) {
	State.writeUint32(KEY, value)
}

function read() {
	return State.readUint32(KEY)
}
`

	worker := NewV8Worker(sdkHandler)
	outputArgs, outputErr, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), contract,
		"write", ArgsToArgumentArray(uint32(1982)))
	require.NoError(t, err)
	require.NoError(t, outputErr)

	outputArgs, outputErr, err = worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), contract,
		"read", ArgsToArgumentArray())
	require.NoError(t, err)
	require.NoError(t, outputErr)

	uin32Value := outputArgs.ArgumentsIterator().NextArguments().Uint32Value()
	require.EqualValues(t, 1982, uin32Value)
}

func TestNewV8Worker_ManipulateStateWithString(t *testing.T) {
	sdkHandler := test.AFakeSdkFor([]byte("signer"), []byte("caller"))

	contract := `
import { State } from "orbs-contract-sdk/v1";
const KEY = new Uint8Array([1, 2, 3])

function write(value) {
	State.writeString(KEY, value)
}

function read() {
	return State.readString(KEY)
}
`

	worker := NewV8Worker(sdkHandler)
	outputArgs, outputErr, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), contract,
		"write", ArgsToArgumentArray("Diamond Dogs"))
	require.NoError(t, err)
	require.NoError(t, outputErr)

	outputArgs, outputErr, err = worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), contract,
		"read", ArgsToArgumentArray())
	require.NoError(t, err)
	require.NoError(t, outputErr)

	uin32Value := outputArgs.ArgumentsIterator().NextArguments().StringValue()
	require.EqualValues(t, "Diamond Dogs", uin32Value)
}
