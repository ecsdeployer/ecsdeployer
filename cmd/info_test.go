package cmd

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInfoCmd(t *testing.T) {
	cmd := newInfoCmd(defaultCmdMetadata())
	cmd.cmd.SetArgs([]string{"-c", "testdata/info_simple.yml"})
	require.NoError(t, cmd.cmd.Execute())
}
