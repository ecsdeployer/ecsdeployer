package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRootCmd(t *testing.T) {
	t.Run("help", func(t *testing.T) {
		result := runCommand("-h")
		require.NoError(t, result.err)
		require.Equal(t, 0, result.exitCode)
		require.Contains(t, result.stdout, "https://ecsdeployer.com/")
	})
	t.Run("version", func(t *testing.T) {
		result := runCommand("-v")
		require.NoError(t, result.err)
		require.Equal(t, fmt.Sprintf("ecsdeployer version %s\n", fakedTestVersionStr), result.stdout)
		require.Equal(t, 0, result.exitCode)
	})
}
