package worker

import "C"
import (
	"fmt"
	"github.com/orbs-network/orbs-contract-sdk/go/context"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/ry/v8worker2"
	"io/ioutil"
)

type wrapper struct {
	sdkHandler context.SdkHandler
	worker *v8worker2.Worker
	callback v8worker2.ReceiveMessageCallback
	value chan interface{}
}

func buildCallback(handler context.SdkHandler, value chan interface{}) v8worker2.ReceiveMessageCallback {
	return func(msg []byte) []byte {
		argArray := protocol.ArgumentArrayReader(msg)

		i := argArray.ArgumentsIterator()
		methodName := i.NextArguments().Uint32Value()
		requestId := i.NextArguments().Uint32Value()

		if methodName == 0 && requestId == 0 {
			value <- ArgsToValue(protocol.ArgumentArrayReader(msg)).Raw()
		} else if methodName == 1 && requestId == 1 {
			addr := handler.SdkAddressGetSignerAddress([]byte("test"), context.PERMISSION_SCOPE_SERVICE)
			return ArgsToArgumentArray(addr).Raw()
		}

		return nil
	}
}

type Worker interface {
	ProcessMethodCall(executionContextId primitives.ExecutionContextId, code string, methodName primitives.MethodName, args *protocol.ArgumentArray) (contractOutputArgs *protocol.ArgumentArray, contractOutputErr error, err error)
}

func (w *wrapper) ProcessMethodCall(executionContextId primitives.ExecutionContextId, code string, methodName primitives.MethodName, args *protocol.ArgumentArray) (contractOutputArgs *protocol.ArgumentArray, contractOutputErr error, err error) {
	wrappedCode, err := WrapWithSDK(code, methodName.String(), nil)
	if err != nil {
		return nil, nil, err
	}

	clientSDK, err := ioutil.ReadFile("../js/arguments.js")
	if err != nil {
		return nil, nil, err
	}

	//textEncoder, err := ioutil.ReadFile("../js/text-encoder-polyfill.js")
	//if err != nil {
	//	return  nil, nil, err
	//}

	w.worker.LoadModule("arguments",
		`const global = {}; export const Arguments = global;` + string(clientSDK), func(moduleName, referrerName string) int {
		println("resolved", moduleName, referrerName)
		return 0
	})

	if err := w.worker.LoadModule(string(executionContextId) + ".js", wrappedCode, func(moduleName, referrerName string) int {
		println("resolved", moduleName, referrerName)
		return 0
	}); err != nil {
		return nil, err, nil
	}

	//println(TypedArgs("message", "methodArguments", args).Raw())

	if err := w.worker.SendBytes(TypedArgs(uint32(999), uint32(666), args).Raw()); err != nil {
		fmt.Println("err!", err)
		return nil, err, nil
	}

	val := (<-w.value).([]byte)
	return protocol.ArgumentArrayReader(val), nil, err
}

func NewV8Worker(sdkHandler context.SdkHandler) Worker {
	// need a buffered channel for return value
	value := make(chan interface{}, 1)
	callback := buildCallback(sdkHandler, value)

	return &wrapper{
		worker: v8worker2.New(callback),
		callback: callback,
		value: value,
		sdkHandler: sdkHandler,
	}
}

var New = NewV8Worker