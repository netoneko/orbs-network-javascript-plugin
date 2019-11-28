package e2e

import (
	"bytes"
	"fmt"
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDeploymentOfJavascriptContract(t *testing.T) {
	skipUnlessE2Enabled(t)

	account, _ := orbs.CreateAccount()
	client := getOrbsClient()

	lt := time.Now()
	PrintTestTime(t, "started", &lt)

	counterStart := uint64(time.Now().UnixNano())
	contractName := fmt.Sprintf("JsTest%d", counterStart)

	PrintTestTime(t, "send deploy - start", &lt)

	DeployContractAndRequireSuccess(t, client, account.PublicKey, account.PrivateKey, contractName, orbs.PROCESSOR_TYPE_JAVASCRIPT,
		[]byte(`
import { State, Address } from "orbs-contract-sdk/v1";
const key = new Uint8Array([1, 2, 3]);

export function _init() {
	State.writeString(key, "Station to Station")
}

export function get() {
	return 100
}

export function getSignerAddress() {
	return Address.getSignerAddress()
}

export function saveName(value) {
	State.writeString(key, value)
}

export function getName() {
	return State.readString(key)
}
`))

	PrintTestTime(t, "send deploy - end", &lt)

	ok := Eventually(EVENTUALLY_DOCKER_E2E_TIMEOUT, func() bool {
		PrintTestTime(t, "run query - start", &lt)
		response, err2 := RunQuery(client, account.PublicKey, contractName, "get")
		PrintTestTime(t, "run query - end", &lt)

		if err2 == nil && response.ExecutionResult == codec.EXECUTION_RESULT_SUCCESS {
			return response.OutputArguments[0].(uint32) == 100
		}
		return false
	})
	require.True(t, ok, "get counter should return counter start")

	signerAddress := account.AddressAsBytes()
	ok = Eventually(EVENTUALLY_DOCKER_E2E_TIMEOUT, func() bool {
		PrintTestTime(t, "run query - start", &lt)
		response, err2 := RunQuery(client, account.PublicKey, contractName, "getSignerAddress")
		PrintTestTime(t, "run query - end", &lt)

		if err2 == nil && response.ExecutionResult == codec.EXECUTION_RESULT_SUCCESS {
			return bytes.Equal(response.OutputArguments[0].([]byte), signerAddress)
		}
		return false
	})
	require.True(t, ok, "getSignerAddress should return signer address")

	ok = Eventually(EVENTUALLY_DOCKER_E2E_TIMEOUT, func() bool {
		response, err := RunQuery(client, account.PublicKey, contractName, "getName")

		if err == nil && response.ExecutionResult == codec.EXECUTION_RESULT_SUCCESS {
			return response.OutputArguments[0].(string) == "Station to Station"
		}
		return false
	})
	require.True(t, ok, "getName should return initial state")

	PrintTestTime(t, "send transaction - start", &lt)
	response, err := SendTransaction(client, account.PublicKey, account.PrivateKey, contractName, "saveName", "Diamond Dogs")
	PrintTestTime(t, "send transaction - end", &lt)

	require.NoError(t, err, "add transaction should not return error")
	require.Equal(t, codec.TRANSACTION_STATUS_COMMITTED, response.TransactionStatus)
	require.Equal(t, codec.EXECUTION_RESULT_SUCCESS, response.ExecutionResult)

	ok = Eventually(EVENTUALLY_DOCKER_E2E_TIMEOUT, func() bool {
		response, err := RunQuery(client, account.PublicKey, contractName, "getName")

		if err == nil && response.ExecutionResult == codec.EXECUTION_RESULT_SUCCESS {
			return response.OutputArguments[0].(string) == "Diamond Dogs"
		}
		return false
	})

	require.True(t, ok, "getName should return name")
	PrintTestTime(t, "done", &lt)

}

func TestDeploymentOfJavascriptContractInteroperableWithGo(t *testing.T) {
	skipUnlessE2Enabled(t)

	account, _ := orbs.CreateAccount()
	client := getOrbsClient()

	lt := time.Now()
	PrintTestTime(t, "started", &lt)

	PrintTestTime(t, "first block committed", &lt)

	counterStart := uint64(time.Now().UnixNano())
	goContractName := fmt.Sprintf("GoTest%d", counterStart)
	jsContractName := fmt.Sprintf("JsTest%d", counterStart)

	PrintTestTime(t, "send deploy - start", &lt)

	DeployContractAndRequireSuccess(t, client, account.PublicKey, account.PrivateKey, goContractName, orbs.PROCESSOR_TYPE_NATIVE,
		[]byte(`
package main

import (
	"github.com/orbs-network/orbs-contract-sdk/go/sdk/v1"
)

var PUBLIC = sdk.Export(getValue, throwPanic)
var SYSTEM = sdk.Export(_init)

func _init() {

}

func getValue() uint64 {
	return uint64(100)
}

func throwPanic() uint64 {
	panic("bang!")
}
`))

	DeployContractAndRequireSuccess(t, client, account.PublicKey, account.PrivateKey, jsContractName, orbs.PROCESSOR_TYPE_JAVASCRIPT,
		[]byte(`
import { Service } from "orbs-contract-sdk/v1";

export function _init() {

}

export function getValue(contractName) {
	return Service.callMethod(contractName, "getValue")
}

export function checkPanic(contractName) {
	return Service.callMethod(contractName, "throwPanic")
}

export function checkNonExistentMethod(contractName) {
	return Service.callMethod(contractName, "methodDoesNotExist")
}
`))

	PrintTestTime(t, "send deploy - end", &lt)

	ok := Eventually(EVENTUALLY_DOCKER_E2E_TIMEOUT, func() bool {
		PrintTestTime(t, "run query - start", &lt)
		response, err2 := RunQuery(client, account.PublicKey, jsContractName, "getValue", goContractName)
		PrintTestTime(t, "run query - end", &lt)

		if err2 == nil && response.ExecutionResult == codec.EXECUTION_RESULT_SUCCESS {
			return response.OutputArguments[0].(uint64) == 100
		}
		return false
	})
	require.True(t, ok, "getValue() should call the go contract and get a result")

	okWithPanic := Eventually(EVENTUALLY_DOCKER_E2E_TIMEOUT, func() bool {
		PrintTestTime(t, "run query - start", &lt)
		response, _ := RunQuery(client, account.PublicKey, jsContractName, "checkPanic", goContractName)
		PrintTestTime(t, "run query - end", &lt)

		if response.ExecutionResult == codec.EXECUTION_RESULT_ERROR_SMART_CONTRACT {
			return response.OutputArguments[0].(string) == "bang!"
		}
		return false
	})
	require.True(t, okWithPanic, "throwPanic() should call the go contract and get an error")

	okWithNonExistentMethod := Eventually(EVENTUALLY_DOCKER_E2E_TIMEOUT, func() bool {
		PrintTestTime(t, "run query - start", &lt)
		response, _ := RunQuery(client, account.PublicKey, jsContractName, "checkNonExistentMethod", goContractName)
		PrintTestTime(t, "run query - end", &lt)

		if response.ExecutionResult == codec.EXECUTION_RESULT_ERROR_SMART_CONTRACT {
			t.Log(response.OutputArguments[0].(string))
			return response.OutputArguments[0].(string) == fmt.Sprintf("method 'methodDoesNotExist' not found on contract '%s'", goContractName)
		}
		return false
	})
	require.True(t, okWithNonExistentMethod, "checkNonExistentMethod() should call the go contract and get an error")
}
