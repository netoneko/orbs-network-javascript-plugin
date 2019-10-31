package worker

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"strings"
	"text/template"
)

var STOP_WORDS = []string{
	"eval",
	"async",
	"V8Worker2",
	"Function",
	"RegExp",
	"WebAssembly",
	"Promise",
}

func SanitizeCode(code string) (string, error) {
	for _, word := range STOP_WORDS {
		if strings.Contains(code, word) {
			return "", errors.New(fmt.Sprintf(`keyword "%s" is forbidden in smart contract code`, word))
		}
	}

	tmpl, err := template.New(`sdk`).Parse(`
Math.random = undefined;

{{.code}}
`)

	if err != nil {
		return "", errors.WithMessage(err, "failed to parse code sanitizer template")
	}

	buf := bytes.NewBufferString("")
	if err = tmpl.Execute(buf, map[string]interface{}{
		"code": code,
	}); err != nil {
		return "", errors.WithMessage(err, "failed to sanitize contract code")
	}

	return buf.String(), nil
}
