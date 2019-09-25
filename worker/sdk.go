package worker

import (
	"bytes"
	"github.com/orbs-network/orbs-spec/types/go/protocol"
	"text/template"
)

func WrapWithSDK(code string, method string, arguments []interface{}) (string, error) {
	tmpl, err := template.New(`sdk`).Parse(`
import { Arguments } from "arguments";
const { argUint32, argUint64, argString, argBytes, argAddress, packedArgumentsEncode, packedArgumentsDecode } = Arguments.Orbs;

/**
  SDK methods start
**/

const Address = {
	GetSignerAddress: () => {
		const response = V8Worker2.send(packedArgumentsEncode([argUint32(100), argUint32(101)]).buffer);
		return packedArgumentsDecode(new Uint8Array(response)).map(a => a.value)[0];
	}
}

const State = {
	WriteBytes: (key, value) => {
		V8Worker2.send(packedArgumentsEncode([argUint32(300), argUint32(301), argBytes(key), argBytes(value)]).buffer);
	},
	ReadBytes: (key) => {
		const response = V8Worker2.send(packedArgumentsEncode([argUint32(300), argUint32(305), argBytes(key)]).buffer);
		return packedArgumentsDecode(new Uint8Array(response)).map(a => a.value)[0];
	}
}

/**
  SDK methods end
**/

function contract(methodCallArguments) {
/** 
  contract code start
**/
{{.code}}
/** 
  contract code end
**/

	return {{.method}}(...methodCallArguments);
}

function serializeReturnValue(val) {
	if (typeof val === "number") {
		return [argUint32(0), argUint32(0), argUint32(val)];
	}

	if (typeof val === "string") {
		return [argUint32(0), argUint32(0), argString(val)];
	}

	if (typeof val === "object") {
		const protoName = val.__proto__.constructor.name;
		if (protoName === "Uint8Array") {
			return [argUint32(0), argUint32(0), argBytes(val)];
		}
	}
}

V8Worker2.recv(function(msg) {
	const [ methodName, requestId, ...methodCallArguments ] = packedArgumentsDecode(new Uint8Array(msg)).map(a => a.value);

	if (methodName === 0) {
		const val = contract(methodCallArguments);
		const payload = packedArgumentsEncode(serializeReturnValue(val));
		V8Worker2.send(payload.buffer);
	}
});
`)

	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")
	if err = tmpl.Execute(buf, map[string]interface{}{
		"code": code,
		"method": method,
		"args": arguments,
	}); err != nil {
		return "", err
	}

	//println(buf.String())

	return buf.String(), nil
}

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

func TypedArgs(messageType uint32, id uint32, args *protocol.ArgumentArray) *protocol.ArgumentArray {
	res := []*protocol.ArgumentBuilder{
		{
			Type: protocol.ARGUMENT_TYPE_UINT_32_VALUE,
			Uint32Value: messageType,
		},
		{
			Type: protocol.ARGUMENT_TYPE_UINT_32_VALUE,
			Uint32Value: id,
		},
	}

	for i := args.ArgumentsIterator(); i.HasNext() ; {
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
