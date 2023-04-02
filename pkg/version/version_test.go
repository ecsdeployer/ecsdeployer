package version_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/version"
	"github.com/stretchr/testify/require"
)

func TestEssential(t *testing.T) {
	require.Equal(t, "master", version.BuildSHA)
	require.Equal(t, "master", version.ShortSHA)
	require.Equal(t, "dev", version.Prerelease)
	require.Equal(t, "9999.0.0", version.DevVersionID)

	require.Equal(t, version.DevVersionID, version.Version)
}

func TestString(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		revertGlobals(t)
		version.Version = "1.2.3"
		version.Prerelease = ""
		require.Equal(t, "1.2.3", version.String())
	})
	t.Run("prerelease", func(t *testing.T) {
		revertGlobals(t)
		version.Version = "1.2.3"
		version.Prerelease = "dev"
		require.Equal(t, "1.2.3-dev", version.String())
	})
}

func TestIsPrerelease(t *testing.T) {

	t.Run("normal", func(t *testing.T) {
		revertGlobals(t)
		version.Prerelease = ""
		require.False(t, version.IsPrelease())
	})
	t.Run("prerelease", func(t *testing.T) {
		revertGlobals(t)
		version.Prerelease = "dev"
		require.True(t, version.IsPrelease())
	})
}

func revertGlobals(t *testing.T) {
	t.Helper()
	oldVersion := version.Version
	oldPrerel := version.Prerelease
	t.Cleanup(func() {
		version.Prerelease = oldPrerel
		version.Version = oldVersion
	})
}
