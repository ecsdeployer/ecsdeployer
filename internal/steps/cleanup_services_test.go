package steps

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestCleanupServicesStep(t *testing.T) {

	tagMatcher := map[string]string{
		"ecsdeployer/project": "dummy",
	}
	serviceArnPrefix := fmt.Sprintf("arn:aws:ecs:%s:%s:service/testcluster/", awsmocker.DefaultRegion, awsmocker.DefaultAccountId)

	t.Run("when no results", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Tagging_GetResources("ecs:service", tagMatcher, []string{}),
		})
		err := CleanupServicesStep(project.Settings.KeepInSync).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("when only relevant services", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Tagging_GetResources("ecs:service", tagMatcher, []string{
				serviceArnPrefix + "dummy-web",
				serviceArnPrefix + "dummy-worker",
			}),
		})
		err := CleanupServicesStep(project.Settings.KeepInSync).Apply(ctx)
		require.NoError(t, err)
	})

	t.Run("when old services", func(t *testing.T) {
		project, ctx := stepTestAwsMocker(t, "testdata/project_advanced.yml", []*awsmocker.MockedEndpoint{
			testutil.Mock_Tagging_GetResources("ecs:service", tagMatcher, []string{
				serviceArnPrefix + "dummy-web",
				serviceArnPrefix + "dummy-worker",
				serviceArnPrefix + "dummy-somethingelse",
			}),
		})
		err := CleanupServicesStep(project.Settings.KeepInSync).Apply(ctx)
		require.NoError(t, err)
	})

}
