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
	emitEvent: () => {
		// FIXME not implemented
	}
}

export const Uint64 = Number; // FIXME later
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

	//println(buf.String())

	return buf.String()
}
