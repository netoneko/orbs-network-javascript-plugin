package test

import (
	"fmt"
	"github.com/netoneko/orbs-network-javascript-plugin/worker"
	"github.com/ry/v8worker2"
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

	constructor := *symbol.(*func(v8worker2.ReceiveMessageCallback) worker.Worker)

	worker := constructor(func(msg []byte) []byte {
		fmt.Println(string(msg))
		require.EqualValues(t, []byte{1, 2, 3, 4, 0, 0, 0, 0}, msg)
		return nil
	})

	err = worker.Load("myScript", `
const buffer = new ArrayBuffer(4*2);
const view = new DataView(buffer);
view.setUint32(0, 1, true);
view.setUint32(1, 2, true);
view.setUint32(2, 3, true);
view.setUint32(3, 4, true);
V8Worker2.send(buffer);
`)
	require.NoError(t, err)
}