package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestManCmd(t *testing.T) {
	silenceLogging(t)

	t.Run("calls correct function", func(t *testing.T) {

		cmd := newManCmd(defaultCmdMetadata()).cmd
		// cmd.Root().SetArgs([]string{"-q"})

		_, _, err := executeCmdAndReturnOutput(cmd)

		require.NoError(t, err)
	})
}
