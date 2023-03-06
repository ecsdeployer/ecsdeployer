package taskdefinition_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil/buildtestutils"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestProxyConfig(t *testing.T) {

	buildtestutils.StartMocker(t)

	t.Run("when not defined", func(t *testing.T) {
		ctx := buildtestutils.LoadProjectConfig(t, "../testdata/baseline.yml")

		pdTask, err := yaml.ParseYAMLString[config.PreDeployTask](`name: testpd1`)
		require.NoError(t, err)

		taskDefinition := genTaskDef(t, ctx, pdTask)

		require.Nil(t, taskDefinition.ProxyConfiguration)

	})

	t.Run("inherit task defaults", func(t *testing.T) {
		ctx := buildtestutils.LoadProjectConfig(t, "../testdata/proxy.yml")

		taskDefinition := genTaskDef(t, ctx, buildtestutils.GetPredeployTask(ctx.Project, "pd-def"))

		require.NotNil(t, taskDefinition.ProxyConfiguration)
		proxy := taskDefinition.ProxyConfiguration
		require.EqualValues(t, "APPMESH", proxy.Type, "Type")
		require.EqualValues(t, "envoy", *proxy.ContainerName, "ContainerName")

		props := buildtestutils.KVListToMap(proxy.Properties, buildtestutils.KVListToMap_KVP)
		require.Equal(t, "1234", props["AppPorts"])
		require.Equal(t, "dummy", props["Farts"])
	})

	t.Run("task override", func(t *testing.T) {
		ctx := buildtestutils.LoadProjectConfig(t, "../testdata/proxy.yml")

		taskDefinition := genTaskDef(t, ctx, buildtestutils.GetPredeployTask(ctx.Project, "pd-override"))

		require.NotNil(t, taskDefinition.ProxyConfiguration)
		proxy := taskDefinition.ProxyConfiguration
		require.EqualValues(t, "APPMESH", proxy.Type, "Type")
		require.EqualValues(t, "envoy", *proxy.ContainerName, "ContainerName")

		props := buildtestutils.KVListToMap(proxy.Properties, buildtestutils.KVListToMap_KVP)
		require.Equal(t, "5678", props["AppPorts"])
		require.NotContains(t, props, "Farts")
	})

	t.Run("task override disabled", func(t *testing.T) {
		ctx := buildtestutils.LoadProjectConfig(t, "../testdata/proxy.yml")

		taskDefinition := genTaskDef(t, ctx, buildtestutils.GetPredeployTask(ctx.Project, "pd-no"))

		require.Nil(t, taskDefinition.ProxyConfiguration)
	})

	t.Run("proper values", func(t *testing.T) {
		ctx := buildtestutils.LoadProjectConfig(t, "../testdata/proxy.yml")

		taskDefinition := genTaskDef(t, ctx, buildtestutils.GetPredeployTask(ctx.Project, "pd-values"))

		require.NotNil(t, taskDefinition.ProxyConfiguration)
		proxy := taskDefinition.ProxyConfiguration
		require.EqualValues(t, "WRONG", proxy.Type, "Type")
		require.EqualValues(t, "blah", *proxy.ContainerName, "Type")

		props := buildtestutils.KVListToMap(proxy.Properties, buildtestutils.KVListToMap_KVP)

		require.Contains(t, props, "SomeInt")
	})
}
