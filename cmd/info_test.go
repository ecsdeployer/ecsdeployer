package cmd

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestInfoCmd(t *testing.T) {
	silenceLogging(t)
	testutil.StartMocker(t, nil)

	cmd := newRootCmd("testing", func(i int) {}).cmd
	cmd.SetArgs([]string{"info", "-c", "testdata/info_simple.yml", "--trace"})

	_, _, err := executeCmdAndReturnOutput(cmd)

	require.NoError(t, err)
}
