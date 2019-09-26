package pack

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
)

func Pack(source string, target string, targetPkg string, targetName string) error {
	clientSDK, err := ioutil.ReadFile(source)
	if err != nil {
		return err
	}

	packedFile := []byte(fmt.Sprintf(`
package %s

import "encoding/base64"

var %s_DATA, _ = base64.StdEncoding.DecodeString("%s")

func %s() []byte {
	return %s_DATA
}
`,
	targetPkg,
	targetName, base64.StdEncoding.EncodeToString(clientSDK),
	targetName, targetName))

	return ioutil.WriteFile(target, packedFile, 0644)
}
