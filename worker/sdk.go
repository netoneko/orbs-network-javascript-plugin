package worker

import (
	"bytes"
	"fmt"
	"text/template"
)

func DefineSDK() string {
	tmpl, err := template.New(`sdk`).Parse(`
import { Arguments } from "arguments";
const { argUint32, argUint64, argString, argBytes, argAddress, packedArgumentsEncode, packedArgumentsDecode } = Arguments.Orbs;

function protoEquals(val, f) {
	return val.__proto__.constructor == f;
}

function isUint8Array(val) {
	return protoEquals(val, Uint8Array)
}

function isError(val) {
	return protoEquals(val, Error) || protoEquals(val, ReferenceError) || protoEquals(val, TypeError);
}

export function toArgument(val) {
	switch(typeof val) {
		case "object":
			if (isUint8Array(val)) {
				return argBytes(val);
			}

			if (isError(val)) {
				return argString(val.message);
			}
		case "number":
			return argUint32(val);
		case "bigint":
			return argUint64(val);
		case "string":
			return argString(val);
		default:
			throw new Error('failed to convert value "' + val + '" to any argument type');
	}
}

export const Types = {
	protoEquals,
	isError,
	isUint8Array,
	toArgument
}

export const Address = {
	getSignerAddress: {{.sdkMethodGetSignerAddress}},
	getCallerAddress: {{.sdkMethodGetCallerAddress}},
	validateAddress: () => {
    	// FIXME address validation is not part of the SDK handler
	}
}

export const State = {
	writeBytes: {{.sdkMethodWriteBytes}},
	readBytes: {{.sdkMethodReadBytes}},
	writeUint32: {{.sdkMethodWriteUint32}},
	readUint32: {{.sdkMethodReadUint32}},
	writeUint64: {{.sdkMethodWriteUint64}},
	readUint64: {{.sdkMethodReadUint64}},
	writeString: {{.sdkMethodWriteString}},
	readString: {{.sdkMethodReadString}},
	clear: {{.sdkMethodClear}},
}

export const Events = {
	emitEvent: function(eventValidator, ...params) {
		(function(V8Worker2) { // safeguard from injections
			eventValidator(...params);
		})();
		const name = eventValidator.name;
		const serializedParams = (params || []).map(toArgument);
		V8Worker2.send(packedArgumentsEncode([argUint32(400), argUint32(401), argString(name), ...serializedParams]).buffer);
	}
}

export const Uint64 = BigInt;
export const Uint32 = Number;

export const Verify = {
	bytes: () => {
		// FIXME not implemented
	},
	uint32: () => {
		// FIXME not implemented
	},
	uint64: () => {
		// FIXME not implemented
	},
}
`)

	if err != nil {
		panic(fmt.Sprintf("failed to parse SDK bindings template: %s", err))
	}

	buf := bytes.NewBufferString("")
	if err = tmpl.Execute(buf, getSDKSettings()); err != nil {
		panic(fmt.Sprintf("failed to generate SDK bindings: %s", err))
	}

	return buf.String()
}
