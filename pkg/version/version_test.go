package version_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/version"
	"github.com/stretchr/testify/require"
)

func TestEssential(t *testing.T) {
	require.Equal(t, "master", version.BuildSHA)
	require.Equal(t, "development", version.DevVersionID)

	require.Equal(t, version.DevVersionID, version.Version)
}
