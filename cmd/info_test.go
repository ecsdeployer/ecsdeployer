package cmd

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestInfoCmd(t *testing.T) {
	silenceLogging(t)
	testutil.StartMocker(t, nil)

	cmd := newInfoCmd(defaultCmdMetadata()).cmd
	cmd.SetArgs([]string{"-c", "testdata/info_simple.yml"})

	_, _, err := executeCmdAndReturnOutput(cmd)

	require.NoError(t, err)
}
