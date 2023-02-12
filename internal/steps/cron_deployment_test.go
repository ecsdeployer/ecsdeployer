package steps

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCronDeploymentStep(t *testing.T) {
	t.Run("using old eventbridge", func(t *testing.T) {
		project, _ := stepTestAwsMocker(t, "testdata/project_advanced.yml", nil)

		project.Settings.CronUsesEventing = true

		step := CronDeploymentStep(project)
		require.True(t, step.ParallelDeps)
		require.Len(t, step.Dependencies, 1)
		require.Equal(t, "Cronjob", step.Dependencies[0].Label)
	})
}
