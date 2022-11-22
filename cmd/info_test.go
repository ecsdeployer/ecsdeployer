package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestInfoCmd(t *testing.T) {
	silenceLogging(t)
	awsmocker.Start(t, nil)

	cmd := newInfoCmd(defaultCmdMetadata()).cmd
	cmd.SetArgs([]string{"-c", "testdata/info_simple.yml"})

	_, _, err := executeCmdAndReturnOutput(cmd)

	require.NoError(t, err)
}
