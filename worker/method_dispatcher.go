package worker

import (
	"encoding/binary"
	"github.com/orbs-network/orbs-contract-sdk/go/context"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
)

const (
	SDK_OBJECT_ADDRESS              = uint32(100)
	SDK_METHOD_GET_SIGNER_ADDRESS   = uint32(101)
	SDK_METHOD_GET_CALLER_ADDRESS   = uint32(102)
	SDK_METHOD_GET_OWN_ADDRESS      = uint32(103)
	SDK_METHOD_GET_CONTRACT_ADDRESS = uint32(104)

	SDK_OBJECT_ENV                  = uint32(200)
	SDK_METHOD_GET_BLOCK_HEIGHT     = uint32(201)
	SDK_METHOD_GET_BLOCK_TIMESTAMP  = uint32(202)
	SDK_METHOD_GET_VIRTUAL_CHAIN_ID = uint32(203)

	SDK_OBJECT_STATE        = uint32(300)
	SDK_METHOD_WRITE_BYTES  = uint32(301)
	SDK_METHOD_WRITE_STRING = uint32(302)
	SDK_METHOD_WRITE_UINT32 = uint32(303)
	SDK_METHOD_WRITE_UINT64 = uint32(304)
	SDK_METHOD_READ_BYTES   = uint32(305)
	SDK_METHOD_READ_STRING  = uint32(306)
	SDK_METHOD_READ_UINT32  = uint32(307)
	SDK_METHOD_READ_UINT64  = uint32(308)
	SDK_METHOD_CLEAR        = uint64(309)
)

type SDKMethodDispatcher interface {
	Dispatch(ctx context.ContextId, permissionScope context.PermissionScope, args *protocol.ArgumentArray) *protocol.ArgumentArray
}

type sdkMethodDispatcher struct {
	handler context.SdkHandler
}

func NewMethodDispatcher(handler context.SdkHandler) SDKMethodDispatcher {
	return &sdkMethodDispatcher{
		handler: handler,
	}
}

func (dispatcher *sdkMethodDispatcher) Dispatch(ctx context.ContextId, permissionScope context.PermissionScope, args *protocol.ArgumentArray) *protocol.ArgumentArray {
	iterator := args.ArgumentsIterator()
	object := iterator.NextArguments().Uint32Value()
	method := iterator.NextArguments().Uint32Value()

	var results []interface{}

	switch object {
	case SDK_OBJECT_ADDRESS:
		switch method {
		case SDK_METHOD_GET_SIGNER_ADDRESS:
			results = append(results, dispatcher.handler.SdkAddressGetSignerAddress(ctx, permissionScope))
		case SDK_METHOD_GET_CALLER_ADDRESS:
			results = append(results, dispatcher.handler.SdkAddressGetCallerAddress(ctx, permissionScope))
		case SDK_METHOD_GET_OWN_ADDRESS:
			results = append(results, dispatcher.handler.SdkAddressGetOwnAddress(ctx, permissionScope))
		case SDK_METHOD_GET_CONTRACT_ADDRESS:
			contractName := iterator.NextArguments().StringValue()
			results = append(results, dispatcher.handler.SdkAddressGetContractAddress(ctx, permissionScope, contractName))
		}
	case SDK_OBJECT_ENV:
		switch method {
		case SDK_METHOD_GET_BLOCK_HEIGHT:
			results = append(results, dispatcher.handler.SdkEnvGetBlockHeight(ctx, permissionScope))
		case SDK_METHOD_GET_BLOCK_TIMESTAMP:
			results = append(results, dispatcher.handler.SdkEnvGetBlockTimestamp(ctx, permissionScope))
		case SDK_METHOD_GET_VIRTUAL_CHAIN_ID:
			results = append(results, dispatcher.handler.SdkEnvGetVirtualChainId(ctx, permissionScope))
		}
	case SDK_OBJECT_STATE:
		key := iterator.NextArguments().BytesValue()

		switch method {
		case SDK_METHOD_WRITE_BYTES:
			value := iterator.NextArguments().BytesValue()
			dispatcher.handler.SdkStateWriteBytes(ctx, permissionScope, key, value)
		case SDK_METHOD_READ_BYTES:
			value := dispatcher.handler.SdkStateReadBytes(ctx, permissionScope, key)
			results = append(results, value)
		case SDK_METHOD_WRITE_STRING:
			value := iterator.NextArguments().StringValue()
			dispatcher.handler.SdkStateWriteBytes(ctx, permissionScope, key, []byte(value))
		case SDK_METHOD_READ_STRING:
			value := dispatcher.handler.SdkStateReadBytes(ctx, permissionScope, key)
			results = append(results, string(value))
		case SDK_METHOD_WRITE_UINT32:
			value := make([]byte, 4)
			binary.LittleEndian.PutUint32(value, iterator.NextArguments().Uint32Value())
			dispatcher.handler.SdkStateWriteBytes(ctx, permissionScope, key, value)
		case SDK_METHOD_READ_UINT32:
			value := dispatcher.handler.SdkStateReadBytes(ctx, permissionScope, key)
			if len(value) < 4 {
				results = append(results, 0)
			} else {
				results = append(results, binary.LittleEndian.Uint32(value))
			}
		case SDK_METHOD_WRITE_UINT64:
			value := make([]byte, 8)
			binary.LittleEndian.PutUint64(value, iterator.NextArguments().Uint64Value())
			dispatcher.handler.SdkStateWriteBytes(ctx, permissionScope, key, value)
		case SDK_METHOD_READ_UINT64:
			value := dispatcher.handler.SdkStateReadBytes(ctx, permissionScope, key)
			if len(value) < 4 {
				results = append(results, 0)
			} else {
				results = append(results, binary.LittleEndian.Uint64(value))
			}
		}
	}

	return ArgsToArgumentArray(results...)
}
