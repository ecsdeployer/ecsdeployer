package taskdefinition_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestProxyConfig(t *testing.T) {

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

		require.Nil(t, taskDefinition.ProxyConfiguration)

	})

	t.Run("inherit task defaults", func(t *testing.T) {
		ctx := loadProjectConfig(t, "proxy.yml")

		taskDefinition := genTaskDef(t, ctx, getPredeployTask(ctx.Project, "pd-def"))

		require.NotNil(t, taskDefinition.ProxyConfiguration)
		proxy := taskDefinition.ProxyConfiguration
		require.EqualValues(t, "APPMESH", proxy.Type, "Type")
		require.EqualValues(t, "envoy", *proxy.ContainerName, "ContainerName")

		props := kvListToMap(proxy.Properties, kvListToMap_KVP)
		require.Equal(t, "1234", props["AppPorts"])
		require.Equal(t, "dummy", props["Farts"])
	})

	t.Run("task override", func(t *testing.T) {
		ctx := loadProjectConfig(t, "proxy.yml")

		taskDefinition := genTaskDef(t, ctx, getPredeployTask(ctx.Project, "pd-override"))

		require.NotNil(t, taskDefinition.ProxyConfiguration)
		proxy := taskDefinition.ProxyConfiguration
		require.EqualValues(t, "APPMESH", proxy.Type, "Type")
		require.EqualValues(t, "envoy", *proxy.ContainerName, "ContainerName")

		props := kvListToMap(proxy.Properties, kvListToMap_KVP)
		require.Equal(t, "5678", props["AppPorts"])
		require.NotContains(t, props, "Farts")
	})

	t.Run("task override disabled", func(t *testing.T) {
		ctx := loadProjectConfig(t, "proxy.yml")

		taskDefinition := genTaskDef(t, ctx, getPredeployTask(ctx.Project, "pd-no"))

		require.Nil(t, taskDefinition.ProxyConfiguration)
	})

	t.Run("proper values", func(t *testing.T) {
		ctx := loadProjectConfig(t, "proxy.yml")

		taskDefinition := genTaskDef(t, ctx, getPredeployTask(ctx.Project, "pd-values"))

		require.NotNil(t, taskDefinition.ProxyConfiguration)
		proxy := taskDefinition.ProxyConfiguration
		require.EqualValues(t, "WRONG", proxy.Type, "Type")
		require.EqualValues(t, "blah", *proxy.ContainerName, "Type")

		props := kvListToMap(proxy.Properties, kvListToMap_KVP)

		require.Contains(t, props, "SomeInt")
	})
}
