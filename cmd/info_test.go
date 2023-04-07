package cmd

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestInfoCmd(t *testing.T) {
	silenceLogging(t)
	testutil.StartMocker(t, nil)

	result := runCommand(t, "info", "-c", "testdata/info_simple.yml", "--trace")
	require.NoError(t, result.err)
}
