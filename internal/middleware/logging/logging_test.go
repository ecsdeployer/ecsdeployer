package logging

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	log "github.com/caarlos0/log"
	"github.com/stretchr/testify/require"
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
