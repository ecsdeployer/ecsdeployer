package pipeline_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/pipeline"
	"github.com/stretchr/testify/require"
)

func TestGlobals(t *testing.T) {
	require.IsType(t, []pipeline.Stepper{}, pipeline.CleanupPipeline)
	require.IsType(t, []pipeline.Stepper{}, pipeline.DeploymentPipeline)
}
