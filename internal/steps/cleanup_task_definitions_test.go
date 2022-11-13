package steps

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestCleanupTaskDefinitionsStep(t *testing.T) {

	tagMatcher := map[string]string{
		"ecsdeployer/project": "dummy",
	}
	taskDefPrefix := fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/", awsmocker.DefaultRegion, awsmocker.DefaultAccountId)

	t.Run("when no results", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Tagging_GetResources("ecs:task-definition", tagMatcher, []string{}),
		})
		err := CleanupTaskDefinitionsStep(project.Settings.KeepInSync).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("when only relevant task defs", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Tagging_GetResources("ecs:task-definition", tagMatcher, []string{
				taskDefPrefix + "dummy-web:1",
				taskDefPrefix + "dummy-worker:1",
			}),
		})
		err := CleanupTaskDefinitionsStep(project.Settings.KeepInSync).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("when task defs for old entries", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Tagging_GetResources("ecs:task-definition", tagMatcher, []string{
				taskDefPrefix + "dummy-web:1",
				taskDefPrefix + "dummy-worker:1",
				taskDefPrefix + "dummy-somethingelse:1",
			}),
			testutil.Mock_ECS_DeregisterTaskDefinition("dummy-somethingelse", 1),
		})
		err := CleanupTaskDefinitionsStep(project.Settings.KeepInSync).Apply(ctx)
		require.NoError(t, err)
	})

}
