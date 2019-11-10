package worker

import (
	"github.com/orbs-network/orbs-contract-sdk/go/context"
	"github.com/orbs-network/orbs-network-javascript-plugin/test"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMethodDispatcher(t *testing.T) {
	handler := test.AFakeSdkFor([]byte("signer"), []byte("caller"))
	dispatcher := NewMethodDispatcher(handler)

	packedSignerAddress, err := dispatcher.Dispatch(context.ContextId("test"), context.PERMISSION_SCOPE_SERVICE,
		ArgsToArgumentArray(SDK_OBJECT_ADDRESS, SDK_METHOD_GET_SIGNER_ADDRESS))
	require.NoError(t, err)
	signerAddress := packedSignerAddress.ArgumentsIterator().NextArguments().BytesValue()
	require.EqualValues(t, []byte("signer"), signerAddress)

	handler.MockEnvBlockHeight(1221)

	packedBlockHeight, err := dispatcher.Dispatch(context.ContextId("test"), context.PERMISSION_SCOPE_SERVICE,
		ArgsToArgumentArray(SDK_OBJECT_ENV, SDK_METHOD_GET_BLOCK_HEIGHT))
	require.NoError(t, err)
	blockHeight := packedBlockHeight.ArgumentsIterator().NextArguments().Uint64Value()
	require.EqualValues(t, 1221, blockHeight)
}

func TestMethodDispatcherWithState(t *testing.T) {
	handler := test.AFakeSdkFor([]byte("signer"), []byte("caller"))
	dispatcher := NewMethodDispatcher(handler)

	dispatcher.Dispatch(context.ContextId("test"), context.PERMISSION_SCOPE_SERVICE,
		ArgsToArgumentArray(SDK_OBJECT_STATE, SDK_METHOD_WRITE_BYTES, []byte("album"), []byte("Diamond Dogs")))

	packedStateValue, err := dispatcher.Dispatch(context.ContextId("test"), context.PERMISSION_SCOPE_SERVICE,
		ArgsToArgumentArray(SDK_OBJECT_STATE, SDK_METHOD_READ_BYTES, []byte("album")))
	require.NoError(t, err)
	album := packedStateValue.ArgumentsIterator().NextArguments().BytesValue()
	require.EqualValues(t, []byte("Diamond Dogs"), album)
}
