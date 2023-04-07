package logging

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/go-log"
)

func TestLogging(t *testing.T) {
	require.NoError(t, Log("foo", func(ctx *config.Context) error {
		return nil
	})(nil))

	require.NoError(t, PadLog("foo", func(ctx *config.Context) error {
		log.Info("a")
		return nil
	})(nil))
}
