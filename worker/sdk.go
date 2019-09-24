package worker

import (
	"bytes"
	"text/template"
)

func WrapWithSDK(code string, method string, arguments []interface{}) (string, error) {
	tmpl, err := template.New(`sdk`).Parse(`
import { Arguments } from "arguments";
const { argUint32, argUint64, argString, argBytes, argAddress, packedArgumentsEncode, packedArgumentsDecode } = Arguments.Orbs;

V8Worker2.print(argUint32);

const val = (function () {
	{{.code}}

	return {{.method}}()
})();

const serializeReturnValue = (val) => {
	const buffer = new ArrayBuffer(4*2);
	const view = new DataView(buffer);

    if (typeof val === "number") {
		view.setUint32(0, val, true);
	}

	return buffer;
}

V8Worker2.send(serializeReturnValue(val));
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