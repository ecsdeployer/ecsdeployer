package cleanuptaskdefinitions

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/testutil/steptestutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestCleanupTaskDefinitionsStep(t *testing.T) {

	t.Run("String", func(t *testing.T) {
		require.Equal(t, "cleaning orphaned task definitions", Step{}.String())
	})

	tagMatcher := map[string]string{
		"ecsdeployer/project": "dummy",
	}
	taskDefPrefix := fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/", awsmocker.DefaultRegion, awsmocker.DefaultAccountId)

	t.Run("when no results", func(t *testing.T) {
		_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Tagging_GetResources("ecs:task-definition", tagMatcher, []string{}),
		})
		err := Step{}.Clean(ctx)
		require.NoError(t, err)
	})

	t.Run("when only relevant task defs", func(t *testing.T) {
		_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Tagging_GetResources("ecs:task-definition", tagMatcher, []string{
				taskDefPrefix + "dummy-web:1",
				taskDefPrefix + "dummy-worker:1",
			}),
		})
		err := Step{}.Clean(ctx)
		require.NoError(t, err)
	})

	t.Run("when task defs for old entries", func(t *testing.T) {
		_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Tagging_GetResources("ecs:task-definition", tagMatcher, []string{
				taskDefPrefix + "dummy-web:1",
				taskDefPrefix + "dummy-worker:1",
				taskDefPrefix + "dummy-somethingelse:1",
			}),
			testutil.Mock_ECS_DeregisterTaskDefinition("dummy-somethingelse", 1),
		})
		err := Step{}.Clean(ctx)
		require.NoError(t, err)
	})

}
