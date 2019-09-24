package test

import (
	"fmt"
	"github.com/netoneko/orbs-network-javascript-plugin/worker"
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

	constructor := *symbol.(*func() worker.Worker)

	worker := constructor()

	outputArgs, outputErr, err := worker.ProcessMethodCall(primitives.ExecutionContextId("myScript"), `
const buffer = new ArrayBuffer(4*2);
const view = new DataView(buffer);
view.setUint32(0, 1, true);
view.setUint32(1, 2, true);
view.setUint32(2, 3, true);
view.setUint32(3, 4, true);
V8Worker2.send(buffer);
`, "sup", nil)
	require.NoError(t, err)
	require.NoError(t, outputErr)

	require.NotNil(t, outputArgs)
}