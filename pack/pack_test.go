package pack

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPackArguments(t *testing.T) {
	err := Pack("../js/arguments.js", "../packed/arguments_packed.go", "packed", "ArgumentsJS")
	require.NoError(t, err)
}
