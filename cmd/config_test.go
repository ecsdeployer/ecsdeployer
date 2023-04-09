package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	silenceLogging(t)
	t.Run("from stdin", func(t *testing.T) {
		origStdin := os.Stdin
		t.Cleanup(func() {
			os.Stdin = origStdin
		})

		const srcFile = "testdata/minimal.yml"

		realProj, err := loadConfig(srcFile)
		require.NoError(t, err)

		f, ferr := os.CreateTemp("", "ecsdestdin.yml")
		require.NoError(t, ferr)
		defer os.Remove(f.Name())
		require.NoError(t, fillStreamWithConfig(t, f, srcFile))
		os.Stdin = f

		proj, err := loadConfig("-")
		require.NoError(t, err)
		require.Equal(t, realProj, proj)
	})
}
