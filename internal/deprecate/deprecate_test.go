package deprecate_test

import (
	"bytes"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/deprecate"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	log "github.com/caarlos0/log"
	"github.com/stretchr/testify/require"
)

func TestDeprecate(t *testing.T) {

	oldLog := log.Log
	t.Cleanup(func() {
		log.Log = oldLog
	})

	var w bytes.Buffer
	log.Log = log.New(&w)

	t.Run("Notice", func(t *testing.T) {
		defer w.Reset()
		ctx := config.New(&config.Project{})
		deprecate.Notice(ctx, "fakething")
		require.Contains(t, w.String(), "DEPRECATED")
	})

}
