package cmd

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestCheckConfig(t *testing.T) {
	closeMock := testutil.MockSimpleStsProxy(t)
	defer closeMock()

	cmd := newCheckCmd(defaultCmdMetadata())
	cmd.cmd.SetArgs([]string{"-c", "testdata/valid.yml"})
	require.NoError(t, cmd.cmd.Execute())
}

func TestCheckConfigThatDoesNotExist(t *testing.T) {
	cmd := newCheckCmd(defaultCmdMetadata())
	cmd.cmd.SetArgs([]string{"-c", "testdata/nope.yml"})
	require.EqualError(t, cmd.cmd.Execute(), "open testdata/nope.yml: no such file or directory")
}

func TestCheckConfigUnmarshalError(t *testing.T) {
	cmd := newCheckCmd(defaultCmdMetadata())
	cmd.cmd.SetArgs([]string{"-c", "testdata/badformat.yml"})
	require.EqualError(t, cmd.cmd.Execute(), "config does not adhere to schema")
}

func TestCheckConfigInvalid(t *testing.T) {
	cmd := newCheckCmd(defaultCmdMetadata())
	cmd.cmd.SetArgs([]string{"-c", "testdata/invalid.yml"})
	require.EqualError(t, cmd.cmd.Execute(), "invalid config: CPU shares provided in an invalid format")
}
