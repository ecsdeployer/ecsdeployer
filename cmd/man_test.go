package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestManCmd(t *testing.T) {
	silenceLogging(t)

	t.Run("misc", func(t *testing.T) {
		cmd := newManCmd().cmd
		require.True(t, cmd.Hidden)
	})

	t.Run("calls correct function", func(t *testing.T) {
		result := runCommand(t, nil, "man")
		require.NoError(t, result.err)
		require.Contains(t, result.stdout, ".TH ECSDEPLOYER")
	})
}
