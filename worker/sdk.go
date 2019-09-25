package worker

import (
	"bytes"
	"text/template"
)

func WrapWithSDK(code string, method string) (string, error) {
	tmpl, err := template.New(`sdk`).Parse(`
import { Arguments } from "arguments";
const { argUint32, argUint64, argString, argBytes, argAddress, packedArgumentsEncode, packedArgumentsDecode } = Arguments.Orbs;

/**
  SDK methods start
**/

const Address = {
	GetSignerAddress: {{.sdkMethodGetSignerAddress}},
	GetCallerAddress: {{.sdkMethodGetCallerAddress}},
}

const State = {
	WriteBytes: {{.sdkMethodWriteBytes}},
	ReadBytes: {{.sdkMethodReadBytes}},
	WriteUint32: {{.sdkMethodWriteUint32}},
	ReadUint32: {{.sdkMethodReadUint32}},
	WriteUint64: {{.sdkMethodWriteUint64}},
	ReadUint64: {{.sdkMethodReadUint64}},
	WriteString: {{.sdkMethodWriteString}},
	ReadString: {{.sdkMethodReadString}},
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
	if err = tmpl.Execute(buf, getSDKCodeSettings(code, method)); err != nil {
		return "", err
	}

	//println(buf.String())

	return buf.String(), nil
}

func proxyWriteOnlyMethodCall(sdkObject uint32, sdkMethod uint32, jsParams string, jsWrappedParams string) string {
	tmpl, err := template.New(`sdkProxyMethodCall`).Parse(`
({{.jsParams}}) => {
	V8Worker2.send(packedArgumentsEncode([argUint32({{.sdkObject}}), argUint32({{.sdkMethod}}), {{.jsWrappedParams}}]).buffer);
}`)

	if err != nil {
		panic(err)
	}

	buf := bytes.NewBufferString("")
	tmpl.Execute(buf, map[string]interface{}{
		"sdkObject": sdkObject,
		"sdkMethod": sdkMethod,
		"jsParams": jsParams,
		"jsWrappedParams": jsWrappedParams,
	})

	return buf.String()
}

// FIXME support contract calls
func proxyReadMethodCall(sdkObject uint32, sdkMethod uint32, jsParams string, jsWrappedParams string) string {
	tmpl, err := template.New(`sdkProxyMethodCall`).Parse(`
({{.jsParams}}) => {
	const response = V8Worker2.send(packedArgumentsEncode([argUint32({{.sdkObject}}), argUint32({{.sdkMethod}}), {{.jsWrappedParams}}]).buffer);
	return packedArgumentsDecode(new Uint8Array(response)).map(a => a.value)[0];
}`)

	if err != nil {
		panic(err)
	}

	buf := bytes.NewBufferString("")
	tmpl.Execute(buf, map[string]interface{}{
		"sdkObject": sdkObject,
		"sdkMethod": sdkMethod,
		"jsParams": jsParams,
		"jsWrappedParams": jsWrappedParams,
	})

	return buf.String()
}

func getSDKCodeSettings(code string, method string) map[string]interface{} {
	return map[string]interface{}{
		"code": code,
		"method": method,
		"sdkMethodGetCallerAddress": proxyReadMethodCall(
			SDK_OBJECT_ADDRESS, SDK_METHOD_GET_CALLER_ADDRESS,
			"", "",
		),
		"sdkMethodGetSignerAddress": proxyReadMethodCall(
			SDK_OBJECT_ADDRESS, SDK_METHOD_GET_SIGNER_ADDRESS,
			"", "",
		),
		"sdkMethodWriteBytes": proxyWriteOnlyMethodCall(
			SDK_OBJECT_STATE, SDK_METHOD_WRITE_BYTES,
			"key, value", "argBytes(key), argBytes(value)",
		),
		"sdkMethodReadBytes": proxyReadMethodCall(
			SDK_OBJECT_STATE, SDK_METHOD_READ_BYTES,
			"key", "argBytes(key)",
		),
		"sdkMethodWriteUint32": proxyWriteOnlyMethodCall(
			SDK_OBJECT_STATE, SDK_METHOD_WRITE_UINT32,
			"key, value", "argBytes(key), argUint32(value)",
		),
		"sdkMethodReadUint32": proxyReadMethodCall(
			SDK_OBJECT_STATE, SDK_METHOD_READ_UINT32,
			"key", "argBytes(key)",
		),
		"sdkMethodWriteUint64": proxyWriteOnlyMethodCall(
			SDK_OBJECT_STATE, SDK_METHOD_WRITE_UINT64,
			"key, value", "argBytes(key), argUint64(value)",
		),
		"sdkMethodReadUint64": proxyReadMethodCall(
			SDK_OBJECT_STATE, SDK_METHOD_READ_UINT64,
			"key", "argBytes(key)",
		),
		"sdkMethodWriteString": proxyWriteOnlyMethodCall(
			SDK_OBJECT_STATE, SDK_METHOD_WRITE_STRING,
			"key, value", "argBytes(key), argString(value)",
		),
		"sdkMethodReadString": proxyReadMethodCall(
			SDK_OBJECT_STATE, SDK_METHOD_READ_STRING,
			"key", "argBytes(key)",
		),
	}
}

