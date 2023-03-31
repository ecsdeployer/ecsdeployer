package deregistertaskdefinitions

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/testutil/steptestutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestDeregisterTaskDefinitionsStep(t *testing.T) {

	t.Run("String", func(t *testing.T) {
		require.Equal(t, "deregistering previous task definitions", Step{}.String())
	})

	t.Run("Run", func(t *testing.T) {

		t.Run("when no results", func(t *testing.T) {
			_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
				testutil.Mock_ECS_ListTaskDefinitionFamilies([]string{}),
			})
			err := Step{}.Run(ctx)
			require.NoError(t, err)
		})

		t.Run("when no old task defs exist", func(t *testing.T) {
			_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
				testutil.Mock_ECS_ListTaskDefinitionFamilies([]string{
					"dummy-web",
					"dummy-worker",
					"dummy-cron1",
					"dummy-cron2",
				}),
				testutil.Mock_ECS_ListTaskDefinitions("dummy-web", []int{1}),
				testutil.Mock_ECS_ListTaskDefinitions("dummy-worker", []int{1}),
				testutil.Mock_ECS_ListTaskDefinitions("dummy-cron1", []int{1}),
				testutil.Mock_ECS_ListTaskDefinitions("dummy-cron2", []int{1}),
			})
			err := Step{}.Run(ctx)
			require.NoError(t, err)
		})

		t.Run("when old defs exist", func(t *testing.T) {
			_, ctx := steptestutil.StepTestAwsMocker(t, "../testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
				testutil.Mock_ECS_ListTaskDefinitionFamilies([]string{
					"dummy-web",
				}),
				testutil.Mock_ECS_ListTaskDefinitions("dummy-web", []int{1, 2, 3, 4}),
				testutil.Mock_ECS_DeregisterTaskDefinition("dummy-web", 1),
				testutil.Mock_ECS_DeregisterTaskDefinition("dummy-web", 2),
				testutil.Mock_ECS_DeregisterTaskDefinition("dummy-web", 3),
			})
			err := Step{}.Run(ctx)
			require.NoError(t, err)
		})
	})

}
