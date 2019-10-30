package worker

import (
	"bytes"
	"text/template"
)

func WrapContract(code string, method string) (string, error) {
	tmpl, err := template.New(`sdk`).Parse(`
import { Arguments } from "arguments";
const { argUint32, argUint64, argString, argBytes, argAddress, packedArgumentsEncode, packedArgumentsDecode } = Arguments.Orbs;

/** 
  contract code start
**/
{{.code}}
/** 
  contract code end
**/

function protoEquals(val, f) {
	return val.__proto__.constructor == f;
}

// FIXME error handling
function serializeReturnValue(val) {
	if (typeof val === "number") {
		return [argUint32(0), argUint32(0), argUint32(val)];
	}

	if (typeof val === "string") {
		return [argUint32(0), argUint32(0), argString(val)];
	}

	if (typeof val === "object") {
		if (protoEquals(val, Uint8Array)) {
			return [argUint32(0), argUint32(0), argBytes(val)];
		}

		if (protoEquals(val, Error)) {
			return [argUint32(0), argUint32(1), argString(val.message)];
		}

		if (protoEquals(val, ReferenceError)) {
			return [argUint32(0), argUint32(1), argString(val.message)];
		}

		if (protoEquals(val, TypeError)) {
			return [argUint32(0), argUint32(1), argString(val.message)];
		}
	}

	if (typeof val === "undefined") {
		return [argUint32(0), argUint32(0)];
	}

	throw new Error("unsupported return value");
}

V8Worker2.recv(function(msg) {
	const [ methodName, requestId, ...methodCallArguments ] = packedArgumentsDecode(new Uint8Array(msg)).map(a => a.value);

	if (methodName === 0) {
		let returnValue;
		try {
			if (typeof {{.method}} === "undefined") {
				throw new Error("method '{{.method}}' not found in contract");
			}

			returnValue = {{.method}}(...methodCallArguments);
		} catch (e) {
			returnValue = e;
			V8Worker2.print(e);
		}

		const payload = packedArgumentsEncode(serializeReturnValue(returnValue));
		V8Worker2.send(payload.buffer);
	}
});
`)

	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")
	if err = tmpl.Execute(buf, getCodeSettings(code, method)); err != nil {
		return "", err
	}

	//println(buf.String())

	return buf.String(), nil
}

func getCodeSettings(code string, method string) map[string]interface{} {
	return map[string]interface{}{
		"code":   code,
		"method": method,
	}
}
