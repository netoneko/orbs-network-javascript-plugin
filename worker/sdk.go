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

const val = (function () {
	{{.code}}

	return {{.method}}()
})();

const serializeReturnValue = (val) => {
    if (typeof val === "number") {
		return [argUint32(val)];
	}

	if (typeof val === "string") {
		return [argString(val)];
	}
}

const payload = packedArgumentsEncode(serializeReturnValue(val));
V8Worker2.send(payload.buffer);
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
