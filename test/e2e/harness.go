package e2e

import (
	"github.com/orbs-network/orbs-client-sdk-go/codec"
	"github.com/orbs-network/orbs-client-sdk-go/orbs"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func DeployContractAndRequireSuccess(t *testing.T, client *orbs.OrbsClient,
	publicKey primitives.EcdsaSecp256K1PublicKey, privateKey primitives.EcdsaSecp256K1PrivateKey, contractName string, processorType orbs.ProcessorType, contractBytes ...[]byte) {
	tx, _, err := client.CreateDeployTransaction(publicKey, privateKey, contractName, processorType, contractBytes...)
	require.NoError(t, err)

	res, err := client.SendTransaction(tx)
	require.NoError(t, err)
	require.EqualValues(t, res.ExecutionResult, codec.EXECUTION_RESULT_SUCCESS)
}

func RunQuery(client *orbs.OrbsClient, senderPublicKey []byte, contractName string, methodName string, args ...interface{}) (response *codec.RunQueryResponse, err error) {
	payload, err := client.CreateQuery(senderPublicKey, contractName, methodName, args...)
	if err != nil {
		return nil, err
	}
	response, err = client.SendQuery(payload)
	return
}

func SendTransaction(client *orbs.OrbsClient,
	publicKey primitives.EcdsaSecp256K1PublicKey, privateKey primitives.EcdsaSecp256K1PrivateKey, contractName string, methodName string, args ...interface{}) (response *codec.SendTransactionResponse, err error) {
	tx, _, err := client.CreateTransaction(publicKey, privateKey, contractName, methodName, args...)
	if err != nil {
		return nil, err
	}

	return client.SendTransaction(tx)
}

func getOrbsClient() *orbs.OrbsClient {
	return orbs.NewClient("http://localhost:8080", 42, codec.NETWORK_TYPE_TEST_NET)
}

func skipUnlessE2Enabled(t *testing.T) {
	if os.Getenv("E2E") != "true" {
		t.Skip("skipping, e2e disabled")
	}
}

func PrintTestTime(t *testing.T, msg string, last *time.Time) {
	t.Logf("%s (+%.3fs)", msg, time.Since(*last).Seconds())
	*last = time.Now()
}

const eventuallyIterations = 20
const EVENTUALLY_DOCKER_E2E_TIMEOUT = 1000 * time.Millisecond

func Eventually(timeout time.Duration, f func() bool) bool {
	for i := 0; i < eventuallyIterations; i++ {
		if testButDontPanic(f) {
			return true
		}
		time.Sleep(timeout / eventuallyIterations)
	}
	return false
}

func testButDontPanic(f func() bool) bool {
	defer func() { recover() }()
	return f()
}
