package cmd

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestInfoCmd(t *testing.T) {
	testutil.DisableLoggingForTest(t)
	testutil.StartMocker(t, nil)

	result := runCommand(t, nil, "info", "-c", "testdata/info_simple.yml", "--trace")
	require.NoError(t, result.err)
}
