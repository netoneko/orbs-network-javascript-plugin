package worker

import (
	"github.com/netoneko/orbs-network-javascript-plugin/test"
	"github.com/orbs-network/orbs-contract-sdk/go/context"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestNewV8Worker_CallMethod(t *testing.T) {
	sdkHandler := test.AFakeSdk()
	worker := NewV8Worker(sdkHandler)
	outputArgs, outputErr, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), `
function hello() {
	return 1
}
`, "hello", ArgsToArgumentArray())
	require.NoError(t, err)
	require.NoError(t, outputErr)
	require.NotNil(t, outputArgs)

	uint32Value := outputArgs.ArgumentsIterator().NextArguments().Uint32Value()
	require.EqualValues(t, 1, uint32Value)
}

func TestNewV8Worker_CallMethodWithArguments(t *testing.T) {
	sdkHandler := test.AFakeSdk()
	worker := NewV8Worker(sdkHandler)
	outputArgs, outputErr, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), `
function hello(a, b) {
	return 1 + a + b
}
`, "hello", ArgsToArgumentArray(uint32(2), uint32(3)))
	require.NoError(t, err)
	require.NoError(t, outputErr)
	require.NotNil(t, outputArgs)

	uint32Value := outputArgs.ArgumentsIterator().NextArguments().Uint32Value()
	require.EqualValues(t, 6, uint32Value)
}

func TestNewV8Worker_CallMethodWithStringArgument(t *testing.T) {
	sdkHandler := test.AFakeSdk()
	worker := NewV8Worker(sdkHandler)
	outputArgs, outputErr, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), `
function hello(a) {
	return "hello, " + a
}
`, "hello", ArgsToArgumentArray("world"))
	require.NoError(t, err)
	require.NoError(t, outputErr)
	require.NotNil(t, outputArgs)

	stringValue := outputArgs.ArgumentsIterator().NextArguments().StringValue()
	require.EqualValues(t, "hello, world", stringValue)
}

func TestNewV8Worker_CallSDKHandlerMethod(t *testing.T) {
	sdkHandler := test.AFakeSdkFor([]byte("signer"), []byte("caller"))

	expectedAddr := sdkHandler.SdkAddressGetSignerAddress([]byte("test"), context.PERMISSION_SCOPE_SERVICE)
	require.EqualValues(t, []byte("signer"), expectedAddr)

	worker := NewV8Worker(sdkHandler)
	outputArgs, outputErr, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), `
function testSignerAddress(a, b, c) {
	const address = Address.GetSignerAddress()
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
const KEY = new Uint8Array([1, 2, 3, 4, 5])

function write(value) {
	State.WriteBytes(KEY, value)
	return 0
}

function read() {
	return State.ReadBytes(KEY)
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
const KEY = new Uint8Array([1, 2, 3])

function write(value) {
	State.WriteUint32(KEY, value)
	return 0
}

function read() {
	return State.ReadUint32(KEY)
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
const KEY = new Uint8Array([1, 2, 3])

function write(value) {
	State.WriteString(KEY, value)
	return 0
}

function read() {
	return State.ReadString(KEY)
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
