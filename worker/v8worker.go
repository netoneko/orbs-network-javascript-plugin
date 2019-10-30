package worker

import (
	"fmt"
	"github.com/orbs-network/orbs-contract-sdk/go/context"
	"github.com/orbs-network/orbs-spec/types/go/primitives"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"github.com/ry/v8worker2"
)

type wrapper struct {
	sdkHandler context.SdkHandler
}

type executionResult struct {
	err   error
	value []byte
}

type Worker interface {
	ProcessMethodCall(executionContextId primitives.ExecutionContextId, code string, methodName primitives.MethodName, args *protocol.ArgumentArray) (contractOutputArgs *protocol.ArgumentArray, contractOutputErr error, err error)
}

var ARGUMENTS_CODE = exportArgumentsJS()
var SDK_CODE = DefineSDK()

func (w *wrapper) ProcessMethodCall(executionContextId primitives.ExecutionContextId, code string, methodName primitives.MethodName, args *protocol.ArgumentArray) (contractOutputArgs *protocol.ArgumentArray, contractOutputErr error, err error) {
	value := make(chan executionResult, 1) // need a buffered channel for return value
	callback := NewMethodDispatcher(w.sdkHandler).GetCallback(value, context.ContextId(executionContextId), context.PERMISSION_SCOPE_SERVICE)
	worker := v8worker2.New(callback)

	//start := time.Now()
	worker.LoadModule("arguments",
		ARGUMENTS_CODE, func(moduleName, referrerName string) int {
			println("resolved", moduleName, referrerName)
			return 0
		})
	//println("Loading arguments library:", time.Since(start).String())

	if err != nil {
		return nil, nil, err
	}
	worker.LoadModule("orbs-contract-sdk/v1", SDK_CODE, func(moduleName, referrerName string) int {
		println("resolved", moduleName, referrerName)
		return 0
	})

	wrappedCode, err := WrapContract(code, methodName.String())
	if err != nil {
		return nil, nil, err
	}
	if err := worker.LoadModule(string(executionContextId)+".js", wrappedCode, func(moduleName, referrerName string) int {
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

	val := <-value
	worker.TerminateExecution()
	return protocol.ArgumentArrayReader(val.value), val.err, err
}

func NewV8Worker(sdkHandler context.SdkHandler) Worker {
	return &wrapper{
		sdkHandler: sdkHandler,
	}
}
