package test

import (
	"encoding/hex"
	"github.com/orbs-network/orbs-network-javascript-plugin/test"
	. "github.com/orbs-network/orbs-network-javascript-plugin/worker"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewV8Worker_BytesToHex(t *testing.T) {
	sdkHandler := test.AFakeSdk()
	contract := `
import { Types } from "orbs-contract-sdk/v1"

export function conAir(bytes) {
	return Types.bytesToHex(bytes)
}`
	worker := newTestWorker(t, sdkHandler, contract)
	question := []byte("Why couldn't you just put the Bunny back in the box?")
	stringValue := worker.callMethodWithoutErrors("conAir", ArgsToArgumentArray(question))

	require.EqualValues(t, hex.EncodeToString(question), stringValue.StringValue())
}

func TestNewV8Worker_HexToBytes(t *testing.T) {
	sdkHandler := test.AFakeSdk()
	contract := `
import { Types } from "orbs-contract-sdk/v1"

export function conAir(hex) {
	return Types.hexToBytes(hex)
}`

	worker := newTestWorker(t, sdkHandler, contract)
	question := []byte("Why couldn't you just put the Bunny back in the box?")
	bytesValue := worker.callMethodWithoutErrors("conAir", ArgsToArgumentArray(hex.EncodeToString(question)))

	require.EqualValues(t, question, bytesValue.BytesValue())
}

func TestNewV8Worker_StringToBytes(t *testing.T) {
	sdkHandler := test.AFakeSdk()
	contract := `
import { Types } from "orbs-contract-sdk/v1"

export function vampiresKiss(val) {
	return Types.stringToBytes(val)
}`
	worker := newTestWorker(t, sdkHandler, contract)
	reply := "You just put it in the right file, according to alphabetical order! Y'know A, B, C, D, E, F, G!"
	bytesValue := worker.callMethodWithoutErrors("vampiresKiss", ArgsToArgumentArray(reply))
	require.EqualValues(t, []byte(reply), bytesValue.BytesValue())
}

func TestNewV8Worker_BytesToString(t *testing.T) {
	sdkHandler := test.AFakeSdk()
	contract := `
import { Types } from "orbs-contract-sdk/v1"

export function vampiresKiss(val) {
	return Types.bytesToString(val)
}`
	worker := newTestWorker(t, sdkHandler, contract)
	reply := "You just put it in the right file, according to alphabetical order! Y'know A, B, C, D, E, F, G!"
	stringValue := worker.callMethodWithoutErrors("vampiresKiss", ArgsToArgumentArray([]byte(reply)))
	require.EqualValues(t, reply, stringValue.StringValue())
}
