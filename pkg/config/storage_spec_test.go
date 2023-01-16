package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestStorageSpec(t *testing.T) {
	spec, err := config.NewStorageSpec(20)
	require.NoError(t, err)
	require.Equal(t, int32(20), spec.Gigabytes())
}
