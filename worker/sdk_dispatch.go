package worker

import (
	"bytes"
	"text/template"
)

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
		"sdkObject":       sdkObject,
		"sdkMethod":       sdkMethod,
		"jsParams":        jsParams,
		"jsWrappedParams": jsWrappedParams,
	})

	return buf.String()
}

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
		"sdkObject":       sdkObject,
		"sdkMethod":       sdkMethod,
		"jsParams":        jsParams,
		"jsWrappedParams": jsWrappedParams,
	})

	return buf.String()
}

func getSDKSettings() map[string]interface{} {
	return map[string]interface{}{
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
