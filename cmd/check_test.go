package cmd

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestCheckConfig(t *testing.T) {
	silenceLogging(t)

	t.Run("happy path", func(t *testing.T) {

		testutil.MockSimpleStsProxy(t)

		cmd := newCheckCmd(defaultCmdMetadata()).cmd
		cmd.SetArgs([]string{"-c", "testdata/valid.yml"})
		_, _, err := executeCmdAndReturnOutput(cmd)
		require.NoError(t, err)
	})

	tables := []struct {
		name          string
		filepath      string
		expectedError string
	}{
		{"DoesNotExist", "testdata/nope.yml", "open testdata/nope.yml: no such file or directory"},
		{"UnmarshalError", "testdata/badformat.yml", "config does not adhere to schema"},
		{"Invalid", "testdata/invalid.yml", "invalid config: CPU shares provided in an invalid format"},
	}

	for _, table := range tables {
		t.Run("failure/"+table.name, func(t *testing.T) {
			cmd := newCheckCmd(defaultCmdMetadata()).cmd
			cmd.SetArgs([]string{"-c", table.filepath})
			_, _, err := executeCmdAndReturnOutput(cmd)
			require.EqualError(t, err, table.expectedError)
		})
	}

}
