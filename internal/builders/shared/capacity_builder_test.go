package shared_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/builders/shared"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestCapacityProviderBuilder(t *testing.T) {

	// TODO: write a test once this is implemented
	require.Panics(t, func() {
		res := &config.Service{}
		shared.CapacityProviderBuilder(res)
	})
}
