package test

import (
	"fmt"
	"github.com/netoneko/orbs-network-javascript-plugin/worker"
	"github.com/orbs-network/orbs-contract-sdk/go/context"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/stretchr/testify/require"
	"os/exec"
	"plugin"
	"testing"
)

type Hello interface {
	Hello() string
}

func Test_Main(t *testing.T) {
	cmd := exec.Command("go", "build", "-buildmode=plugin", "-o", "./main.bin", "..")
	out, err := cmd.CombinedOutput()
	fmt.Println(string(out))
	require.NoError(t, err)

	plug, err := plugin.Open("./main.bin")
	require.NoError(t, err)

	symbol, err := plug.Lookup("Test")
	require.NoError(t, err)

	h := symbol.(Hello)
	require.EqualValues(t, "hello", h.Hello())
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

	outputArgs, outputErr, err := v8Worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), `
function hello(a, b) {
	return 1 + a + b
}
`, "hello", worker.ArgsToArgumentArray(uint32(2), uint32(3)))
	require.NoError(t, err)
	require.NoError(t, outputErr)
	require.NotNil(t, outputArgs)

	uint32Value := outputArgs.ArgumentsIterator().NextArguments().Uint32Value()
	require.EqualValues(t, 6, uint32Value)
}
