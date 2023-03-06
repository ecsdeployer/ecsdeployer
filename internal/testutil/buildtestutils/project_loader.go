package buildtestutils

import (
	"fmt"
	"strings"
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

type optsLoadProjectConfig struct {
	NumSSMVars int
}

type optsLPCFunc func(*optsLoadProjectConfig)

func OptSetNumSSMVars(num int) optsLPCFunc {
	return func(opts *optsLoadProjectConfig) {
		opts.NumSSMVars = num
	}
}

func LoadProjectConfig(t *testing.T, filepath string, optFuncs ...optsLPCFunc) *config.Context {
	t.Helper()

	if !strings.Contains(filepath, "/") {
		filepath = fmt.Sprintf("testdata/%s", filepath)
	}

	ctx, err := config.NewFromYAML(filepath)
	require.NoError(t, err)

	opts := &optsLoadProjectConfig{}
	for _, optFunc := range optFuncs {
		optFunc(opts)
	}

	if opts.NumSSMVars > 0 {
		ctx.Cache.SSMSecrets = make(map[string]config.EnvVar)
		for i := 0; i < opts.NumSSMVars; i++ {
			varKey := fmt.Sprintf("SSM_VAR_%02d", i+1)
			varPath := fmt.Sprintf("/fake/path/secret%02d", i+1)
			ctx.Cache.SSMSecrets[varKey] = config.NewEnvVar(config.EnvVarTypeSSM, varPath)
		}
	}

	return ctx
}
