package worker

import (
	"bytes"
	"github.com/orbs-network/orbs-contract-sdk/go/context"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/pkg/errors"
	"github.com/ry/v8worker2"
	"text/template"
)

const SDK_RETURN_FROM_METHOD = 0
const SDK_RETURN_VALUE = 0
const SDK_RETURN_ERROR = 1

const SDK_GENERIC_EXECUTION_ERROR = "JS contract execution failed"

func sdkDispatchCallback(dispatcher SDKMethodDispatcher, value chan executionResult, ctx context.ContextId, scope context.PermissionScope) v8worker2.ReceiveMessageCallback {
	return func(msg []byte) []byte {
		argArray := protocol.ArgumentArrayReader(msg)

		i := argArray.ArgumentsIterator()
		methodName := i.NextArguments().Uint32Value()
		requestId := i.NextArguments().Uint32Value()

		if methodName == SDK_RETURN_FROM_METHOD && requestId == SDK_RETURN_VALUE {
			value <- executionResult{nil, ArgsToValue(protocol.ArgumentArrayReader(msg)).Raw()}
		} else if methodName == SDK_RETURN_FROM_METHOD && requestId == SDK_RETURN_ERROR {
			value <- executionResult{errors.New(SDK_GENERIC_EXECUTION_ERROR), ArgsToValue(protocol.ArgumentArrayReader(msg)).Raw()}
		} else {
			return dispatcher.Dispatch(ctx, scope, argArray).Raw()
		}

		return nil
	}
}

func proxyWriteOnlyMethodCall(sdkObject uint32, sdkMethod uint32, jsParams string, jsWrappedParams string) string {
	tmpl, err := template.New(`sdkProxyMethodCall`).Parse(`
({{.jsParams}}) => {
	V8Worker2.send(packedArgumentsEncode([argUint32({{.sdkObject}}), argUint32({{.sdkMethod}}), {{.jsWrappedParams}}]).buffer);
}`)

	if err != nil {
		panic(err)
	}

	buf := bytes.NewBufferString("")
	tmpl.Execute(buf, map[string]interface{}{
		"sdkObject":       sdkObject,
		"sdkMethod":       sdkMethod,
		"jsParams":        jsParams,
		"jsWrappedParams": jsWrappedParams,
	})

	return buf.String()
}

func proxyReadMethodCall(sdkObject uint32, sdkMethod uint32, jsParams string, jsWrappedParams string) string {
	tmpl, err := template.New(`sdkProxyMethodCall`).Parse(`
({{.jsParams}}) => {
	const response = V8Worker2.send(packedArgumentsEncode([argUint32({{.sdkObject}}), argUint32({{.sdkMethod}}), {{.jsWrappedParams}}]).buffer);
	return packedArgumentsDecode(new Uint8Array(response)).map(a => a.value)[0];
}`)

	if err != nil {
		panic(err)
	}

	buf := bytes.NewBufferString("")
	tmpl.Execute(buf, map[string]interface{}{
		"sdkObject":       sdkObject,
		"sdkMethod":       sdkMethod,
		"jsParams":        jsParams,
		"jsWrappedParams": jsWrappedParams,
	})

	return buf.String()
}

func getSDKSettings() map[string]interface{} {
	return map[string]interface{}{
		"sdkMethodGetCallerAddress": proxyReadMethodCall(
			SDK_OBJECT_ADDRESS, SDK_METHOD_GET_CALLER_ADDRESS,
			"", "",
		),
		"sdkMethodGetSignerAddress": proxyReadMethodCall(
			SDK_OBJECT_ADDRESS, SDK_METHOD_GET_SIGNER_ADDRESS,
			"", "",
		),
		"sdkMethodWriteBytes": proxyWriteOnlyMethodCall(
			SDK_OBJECT_STATE, SDK_METHOD_WRITE_BYTES,
			"key, value", "argBytes(key), argBytes(value)",
		),
		"sdkMethodReadBytes": proxyReadMethodCall(
			SDK_OBJECT_STATE, SDK_METHOD_READ_BYTES,
			"key", "argBytes(key)",
		),
		"sdkMethodWriteUint32": proxyWriteOnlyMethodCall(
			SDK_OBJECT_STATE, SDK_METHOD_WRITE_UINT32,
			"key, value", "argBytes(key), argUint32(value)",
		),
		"sdkMethodReadUint32": proxyReadMethodCall(
			SDK_OBJECT_STATE, SDK_METHOD_READ_UINT32,
			"key", "argBytes(key)",
		),
		"sdkMethodWriteUint64": proxyWriteOnlyMethodCall(
			SDK_OBJECT_STATE, SDK_METHOD_WRITE_UINT64,
			"key, value", "argBytes(key), argUint64(value)",
		),
		"sdkMethodReadUint64": proxyReadMethodCall(
			SDK_OBJECT_STATE, SDK_METHOD_READ_UINT64,
			"key", "argBytes(key)",
		),
		"sdkMethodWriteString": proxyWriteOnlyMethodCall(
			SDK_OBJECT_STATE, SDK_METHOD_WRITE_STRING,
			"key, value", "argBytes(key), argString(value)",
		),
		"sdkMethodReadString": proxyReadMethodCall(
			SDK_OBJECT_STATE, SDK_METHOD_READ_STRING,
			"key", "argBytes(key)",
		),
	}
}
