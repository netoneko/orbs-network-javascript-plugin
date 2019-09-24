package worker

import "C"
import (
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/ry/v8worker2"
	"io/ioutil"
)

type wrapper struct {
	worker *v8worker2.Worker
	callback v8worker2.ReceiveMessageCallback
	value chan interface{}
}

func buildCallback(value chan interface{}) v8worker2.ReceiveMessageCallback {
	return func(msg []byte) []byte {
		value <- msg
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

	w.worker.LoadModule("arguments", `const global = {}; export const Arguments = global;` + string(clientSDK), func(moduleName, referrerName string) int {
		println("resolved", moduleName, referrerName)
		return 0
	});
	err = w.worker.LoadModule(string(executionContextId) + ".js", wrappedCode, func(moduleName, referrerName string) int {
		println("resolved", moduleName, referrerName)
		return 0
	})
	if err != nil {
		return nil, err, nil
	}

	val := (<-w.value).([]byte)
	return argsToArgumentArray(val), nil, err
}

func argsToArgumentArray(args ...interface{}) *protocol.ArgumentArray {
	res := []*protocol.ArgumentBuilder{}
	for _, arg := range args {
		switch arg.(type) {
		case uint32:
			res = append(res, &protocol.ArgumentBuilder{Type: protocol.ARGUMENT_TYPE_UINT_32_VALUE, Uint32Value: arg.(uint32)})
		case uint64:
			res = append(res, &protocol.ArgumentBuilder{Type: protocol.ARGUMENT_TYPE_UINT_64_VALUE, Uint64Value: arg.(uint64)})
		case string:
			res = append(res, &protocol.ArgumentBuilder{Type: protocol.ARGUMENT_TYPE_STRING_VALUE, StringValue: arg.(string)})
		case []byte:
			res = append(res, &protocol.ArgumentBuilder{Type: protocol.ARGUMENT_TYPE_BYTES_VALUE, BytesValue: arg.([]byte)})
		}
	}
	return (&protocol.ArgumentArrayBuilder{Arguments: res}).Build()
}

func NewV8Worker() Worker {
	// need a buffered channel for return value
	value := make(chan interface{}, 1)
	callback := buildCallback(value)

	return &wrapper{
		worker: v8worker2.New(callback),
		callback: callback,
		value: value,
	}
}

var New = NewV8Worker