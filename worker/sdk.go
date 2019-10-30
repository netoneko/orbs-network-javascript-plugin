package worker

import (
	"bytes"
	"text/template"
)

func DefineSDK() (string, error) {
	tmpl, err := template.New(`sdk`).Parse(`
import { Arguments } from "arguments";
const { argUint32, argUint64, argString, argBytes, argAddress, packedArgumentsEncode, packedArgumentsDecode } = Arguments.Orbs;

/**
  SDK methods start
**/

export const Address = {
	getSignerAddress: {{.sdkMethodGetSignerAddress}},
	getCallerAddress: {{.sdkMethodGetCallerAddress}},
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
}

export const Uint64 = Number; // FIXME later
export const Uint32 = Number;

export const Verify = {
	bytes: () => {
		// FIXME not implemented
	},
	uint64: () => {
		// FIXME not implemented
	},
}

/**
  SDK methods end
**/
`)

	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")
	if err = tmpl.Execute(buf, getSDKSettings()); err != nil {
		return "", err
	}

	//println(buf.String())

	return buf.String(), nil
}
