package cmd

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/commands/mancmd"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestManCmd(t *testing.T) {
	testutil.DisableLoggingForTest(t)

	t.Run("misc", func(t *testing.T) {
		cmd := mancmd.New()
		require.True(t, cmd.Hidden)
	})

	t.Run("calls correct function", func(t *testing.T) {
		result := runCommand(t, nil, "man")
		require.NoError(t, result.err)
		require.Contains(t, result.stdout, ".TH ECSDEPLOYER")
	})
}
