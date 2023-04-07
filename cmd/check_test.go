package cmd

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/stretchr/testify/require"
)

func TestCheckConfig(t *testing.T) {
	silenceLogging(t)
	helpers.IsTestingMode = true

	t.Run("happy path", func(t *testing.T) {
		testutil.MockSimpleStsProxy(t)

		const testFile = "testdata/valid.yml"

		tables := []struct {
			label string
			args  []string
		}{
			{"normal", []string{"check", "-c", testFile}},
			{"quiet", []string{"check", "-c", testFile, "-q"}},
			{"show json", []string{"check", "-c", testFile, "--show"}},
			{"show yaml", []string{"check", "-c", testFile, "--yaml"}},
		}

		for _, table := range tables {
			t.Run(table.label, func(t *testing.T) {

				result := runCommand(t, table.args...)
				require.NoError(t, result.err)
				require.Equal(t, 0, result.exitCode)
			})
		}
	})

	t.Run("failures", func(t *testing.T) {

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
			t.Run(table.name, func(t *testing.T) {
				result := runCommand(t, "check", "-c", table.filepath)
				require.Error(t, result.err)
				require.EqualError(t, result.err, table.expectedError)
			})
		}
	})

}
