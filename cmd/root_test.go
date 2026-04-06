package cmd

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/version"
	"github.com/stretchr/testify/require"
)

func TestRootCmd(t *testing.T) {
	t.Run("help", func(t *testing.T) {
		result := runCommand(t, nil, "-h")
		require.NoError(t, result.err)
		require.Equal(t, 0, result.exitCode)
		require.Contains(t, result.stdout, "https://ecsdeployer.com/")
	})
	t.Run("version", func(t *testing.T) {
		result := runCommand(t, nil, "-v")
		require.NoError(t, result.err)
		require.Equal(t, fmt.Sprintf("ecsdeployer version %s\n", version.String()), result.stdout)
		require.Equal(t, 0, result.exitCode)
	})
}
