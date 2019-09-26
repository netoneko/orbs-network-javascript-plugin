package worker

import "C"
import (
	"fmt"
	"github.com/netoneko/orbs-network-javascript-plugin/packed"
	"github.com/orbs-network/orbs-contract-sdk/go/context"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/ry/v8worker2"
)

type wrapper struct {
	sdkHandler context.SdkHandler
}

func buildCallback(dispatcher SDKMethodDispatcher, value chan interface{}, ctx context.ContextId, scope context.PermissionScope) v8worker2.ReceiveMessageCallback {
	return func(msg []byte) []byte {
		argArray := protocol.ArgumentArrayReader(msg)

		i := argArray.ArgumentsIterator()
		methodName := i.NextArguments().Uint32Value()
		requestId := i.NextArguments().Uint32Value()

		if methodName == 0 && requestId == 0 {
			value <- ArgsToValue(protocol.ArgumentArrayReader(msg)).Raw()
		} else {
			return dispatcher.Dispatch(ctx, scope, argArray).Raw()
		}

		return nil
	}
}

type Worker interface {
	ProcessMethodCall(executionContextId primitives.ExecutionContextId, code string, methodName primitives.MethodName, args *protocol.ArgumentArray) (contractOutputArgs *protocol.ArgumentArray, contractOutputErr error, err error)
}

func (w *wrapper) ProcessMethodCall(executionContextId primitives.ExecutionContextId, code string, methodName primitives.MethodName, args *protocol.ArgumentArray) (contractOutputArgs *protocol.ArgumentArray, contractOutputErr error, err error) {
	value := make(chan interface{}, 1) 	// need a buffered channel for return value
	callback := buildCallback(NewMethodDispatcher(w.sdkHandler), value, context.ContextId(executionContextId), context.PERMISSION_SCOPE_SERVICE)
	worker := v8worker2.New(callback)

	wrappedCode, err := WrapWithSDK(code, methodName.String())
	if err != nil {
		return nil, nil, err
	}

	worker.LoadModule("arguments",
		`const global = {}; export const Arguments = global;` + string(packed.ArgumentsJS()), func(moduleName, referrerName string) int {
		println("resolved", moduleName, referrerName)
		return 0
	})

	if err := worker.LoadModule(string(executionContextId) + ".js", wrappedCode, func(moduleName, referrerName string) int {
		println("resolved", moduleName, referrerName)
		return 0
	}); err != nil {
		return nil, err, nil
	}

	// Could be replaced with a call to get arguments and method name
	if err := worker.SendBytes(TypedArgs(uint32(0), uint32(0), args).Raw()); err != nil {
		fmt.Println("err!", err)
		return nil, err, nil
	}

	val := (<-value).([]byte)
	valCopy := make([]byte, len(val))
	copy(valCopy, val)
	worker.TerminateExecution()
	return protocol.ArgumentArrayReader(val), nil, err
}

func NewV8Worker(sdkHandler context.SdkHandler) Worker {
	return &wrapper{
		sdkHandler: sdkHandler,
	}
}

var New = NewV8Worker