package taskdefinition_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestHealthCheck(t *testing.T) {

	testutil.StartMocker(t, &awsmocker.MockerOptions{
		Mocks: []*awsmocker.MockedEndpoint{
			testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
		},
	})

	t.Run("when not defined", func(t *testing.T) {
		ctx := loadProjectConfig(t, "baseline.yml")

		pdTask, err := yaml.ParseYAMLString[config.PreDeployTask](`name: testpd1`)
		require.NoError(t, err)

		taskDefinition := genTaskDef(t, ctx, pdTask)

		require.Nil(t, taskDefinition.ContainerDefinitions[0].HealthCheck)

	})

	t.Run("inherit task defaults", func(t *testing.T) {
		ctx := loadProjectConfig(t, "healthcheck.yml")

		taskDefinition := genTaskDef(t, ctx, ctx.Project.Services[1])

		require.NotNil(t, taskDefinition.ContainerDefinitions[0].HealthCheck)
		hc := taskDefinition.ContainerDefinitions[0].HealthCheck
		require.EqualValues(t, []string{"CMD", "test", "healthcheck"}, hc.Command, "Command")
	})

	t.Run("task override", func(t *testing.T) {
		ctx := loadProjectConfig(t, "healthcheck.yml")

		taskDefinition := genTaskDef(t, ctx, ctx.Project.Services[0])

		require.NotNil(t, taskDefinition.ContainerDefinitions[0].HealthCheck)
		hc := taskDefinition.ContainerDefinitions[0].HealthCheck
		require.EqualValues(t, []string{"CMD-SHELL", "blah", "yar"}, hc.Command, "Command")
	})

	t.Run("task override disabled", func(t *testing.T) {
		ctx := loadProjectConfig(t, "healthcheck.yml")

		taskDefinition := genTaskDef(t, ctx, ctx.Project.Services[2])

		require.Nil(t, taskDefinition.ContainerDefinitions[0].HealthCheck)
	})

	t.Run("sidecar", func(t *testing.T) {
		ctx := loadProjectConfig(t, "healthcheck.yml")

		taskDefinition := genTaskDef(t, ctx, ctx.Project.Services[0])

		sc, err := getContainer(taskDefinition, "sc1")
		require.NoError(t, err)

		require.NotNil(t, sc.HealthCheck)
		require.EqualValues(t, []string{"CMD", "sc1"}, sc.HealthCheck.Command, "Command")
	})

	t.Run("proper values", func(t *testing.T) {
		ctx := loadProjectConfig(t, "healthcheck.yml")

		taskDefinition := genTaskDef(t, ctx, ctx.Project.PreDeployTasks[0])

		require.NotNil(t, taskDefinition.ContainerDefinitions[0].HealthCheck)
		hc := taskDefinition.ContainerDefinitions[0].HealthCheck
		require.EqualValues(t, []string{"CMD", "test2", "blah"}, hc.Command, "Command")
		require.EqualValues(t, 60, *hc.Interval, "Interval")
		require.EqualValues(t, 2, *hc.Retries, "Retries")
		require.EqualValues(t, 120, *hc.StartPeriod, "StartPeriod")
		require.EqualValues(t, 5, *hc.Timeout, "Timeout")
	})
}
