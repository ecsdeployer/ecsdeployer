package cmdutil_test

import (
	"os"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/util/cmdutil"
	"github.com/stretchr/testify/require"
)

func TestLoadConfig(t *testing.T) {
	testutil.DisableLoggingForTest(t)
	t.Run("from stdin", func(t *testing.T) {
		origStdin := os.Stdin
		t.Cleanup(func() {
			os.Stdin = origStdin
		})

		const srcFile = "../../../cmd/testdata/minimal.yml"

		realProj, err := cmdutil.LoadConfig(srcFile)
		require.NoError(t, err)

		f, ferr := os.CreateTemp("", "ecsdestdin.yml")
		require.NoError(t, ferr)
		defer os.Remove(f.Name())
		require.NoError(t, testutil.FillStreamWithConfig(t, f, srcFile))
		os.Stdin = f

		proj, err := cmdutil.LoadConfig("-")
		require.NoError(t, err)
		require.Equal(t, realProj, proj)
	})
}
