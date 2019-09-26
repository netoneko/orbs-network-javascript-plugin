package test

import (
	"fmt"
	"github.com/orbs-network/orbs-network-javascript-plugin/pack"
	"github.com/orbs-network/orbs-network-javascript-plugin/worker"
	"github.com/orbs-network/orbs-contract-sdk/go/context"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/stretchr/testify/require"
	"os/exec"
	"plugin"
	"testing"
)

func TestPackArguments(t *testing.T) {
	err := pack.Pack("../js/arguments.js", "../packed/arguments_packed.go", "packed", "ArgumentsJS")
	require.NoError(t, err)
}

func Test_V8Worker(t *testing.T) {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", "./main.bin", "..")
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	require.NoError(t, err)

	plug, err := plugin.Open("./main.bin")
	require.NoError(t, err)

	symbol, err := plug.Lookup("New")
	require.NoError(t, err)

	constructor := *symbol.(*func(context.SdkHandler) worker.Worker)

	fakeSDK := AFakeSdk()
	v8Worker := constructor(fakeSDK)

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

	outputArgs, outputErr, err := v8Worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), contract,
		"write", worker.ArgsToArgumentArray([]byte("Diamond Dogs")))
	require.NoError(t, err)
	require.NoError(t, outputErr)

	outputArgs, outputErr, err = v8Worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), contract,
		"read", worker.ArgsToArgumentArray())
	require.NoError(t, err)
	require.NoError(t, outputErr)

	bytesValue := outputArgs.ArgumentsIterator().NextArguments().BytesValue()
	require.EqualValues(t, []byte("Diamond Dogs"), bytesValue)
}
