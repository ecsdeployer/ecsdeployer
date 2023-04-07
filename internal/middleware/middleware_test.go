package middleware_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/middleware"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestGlobals(t *testing.T) {
	dunno := func(x middleware.Action) error { return x(nil) }
	exFunc := func(ctx *config.Context) error { return nil }
	require.Nil(t, dunno(exFunc))
}
