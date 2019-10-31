package worker

import (
	"github.com/orbs-network/orbs-network-javascript-plugin/packed"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
)

func ArgsToArgumentArray(args ...interface{}) *protocol.ArgumentArray {
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

func ArgumentArrayToArgs(ArgumentArray *protocol.ArgumentArray) []interface{} {
	res := []interface{}{}
	for i := ArgumentArray.ArgumentsIterator(); i.HasNext(); {
		Argument := i.NextArguments()
		switch Argument.Type() {
		case protocol.ARGUMENT_TYPE_UINT_32_VALUE:
			res = append(res, Argument.Uint32Value())
		case protocol.ARGUMENT_TYPE_UINT_64_VALUE:
			res = append(res, Argument.Uint64Value())
		case protocol.ARGUMENT_TYPE_STRING_VALUE:
			res = append(res, Argument.StringValue())
		case protocol.ARGUMENT_TYPE_BYTES_VALUE:
			res = append(res, Argument.BytesValue())
		}
	}
	return res
}

func TypedArgs(messageType uint32, id uint32, args *protocol.ArgumentArray) *protocol.ArgumentArray {
	res := []*protocol.ArgumentBuilder{
		{
			Type:        protocol.ARGUMENT_TYPE_UINT_32_VALUE,
			Uint32Value: messageType,
		},
		{
			Type:        protocol.ARGUMENT_TYPE_UINT_32_VALUE,
			Uint32Value: id,
		},
	}

	for i := args.ArgumentsIterator(); i.HasNext(); {
		res = append(res, protocol.ArgumentBuilderFromRaw(i.NextArguments().Raw()))
	}

	return (&protocol.ArgumentArrayBuilder{Arguments: res}).Build()
}

func ArgsToValue(args *protocol.ArgumentArray) *protocol.ArgumentArray {
	res := []*protocol.ArgumentBuilder{}

	i := args.ArgumentsIterator()

	// skip 2 steps removing type info
	i.NextArguments()
	i.NextArguments()

	for i.HasNext() {
		res = append(res, protocol.ArgumentBuilderFromRaw(i.NextArguments().Raw()))
	}

	return (&protocol.ArgumentArrayBuilder{Arguments: res}).Build()
}

func exportArgumentsJS() string {
	return `const global = {}; export const Arguments = global;` + string(packed.ArgumentsJS())
}
