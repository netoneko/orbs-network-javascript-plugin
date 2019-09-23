package worker

import "C"
import (
	"github.com/ry/v8worker2"
)

type ReceiveMessageCallback v8worker2.ReceiveMessageCallback

type ModuleResolverCallback v8worker2.ModuleResolverCallback

type Worker interface {
	Dispose()

	Load(scriptName string, code string) error
	LoadModule(scriptName string, code string, resolve v8worker2.ModuleResolverCallback) error
	SendBytes(msg []byte) error
	TerminateExecution()
}

func NewV8Worker(cb v8worker2.ReceiveMessageCallback) Worker {
	return v8worker2.New(cb)
}

//func ResolveModule(moduleSpecifier *C.char, referrerSpecifier *C.char, resolverToken int) C.int {
//	return v8worker2.ResolveModule(moduleSpecifier, referrerSpecifier, resolverToken)
//}

func SetFlags(args []string) []string {
	return v8worker2.SetFlags(args)
}

func Version() string {
	return v8worker2.Version()
}

var New = NewV8Worker