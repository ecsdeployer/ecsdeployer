package taskdefinition_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil/buildtestutils"
	"github.com/stretchr/testify/require"
)

func TestSidecars(t *testing.T) {

	buildtestutils.StartMocker(t)

	t.Run("inherit_env", func(t *testing.T) {
		ctx := buildtestutils.LoadProjectConfig(t, "../testdata/everything.yml")

		task := buildtestutils.GetPredeployTask(ctx.Project, "pd-sc-inherit")
		taskDefinition := genTaskDef(t, ctx, task)

		t.Run("enabled", func(t *testing.T) {
			sc1, err := buildtestutils.GetContainer(taskDefinition, "sc1")
			require.NoError(t, err)

			scEnv := buildtestutils.KVListToMap(sc1.Environment, buildtestutils.KVListToMap_KVP)
			require.Contains(t, scEnv, "SC_TEST_VAR")
			require.Equal(t, "blah", scEnv["SC_TEST_VAR"])
		})

		t.Run("disabled", func(t *testing.T) {
			sc1, err := buildtestutils.GetContainer(taskDefinition, "scno")
			require.NoError(t, err)

			scEnv := buildtestutils.KVListToMap(sc1.Environment, buildtestutils.KVListToMap_KVP)
			require.NotContains(t, scEnv, "SC_TEST_VAR")
		})

		t.Run("default", func(t *testing.T) {
			sc1, err := buildtestutils.GetContainer(taskDefinition, "scdef")
			require.NoError(t, err)

			scEnv := buildtestutils.KVListToMap(sc1.Environment, buildtestutils.KVListToMap_KVP)
			require.NotContains(t, scEnv, "SC_TEST_VAR")
		})

		t.Run("extra_vars_added", func(t *testing.T) {
			sc, err := buildtestutils.GetContainer(taskDefinition, "scadd")
			require.NoError(t, err)

			scEnv := buildtestutils.KVListToMap(sc.Environment, buildtestutils.KVListToMap_KVP)
			require.Contains(t, scEnv, "EXTRA_VAR")
			require.Contains(t, scEnv, "TPL_MULTI_TEST")
			require.Contains(t, scEnv, "SC_TEST_VAR")
			require.Equal(t, "someval", scEnv["EXTRA_VAR"])
			require.Equal(t, "overridden", scEnv["SC_TEST_VAR"])
		})

		t.Run("env_var_no_inherit", func(t *testing.T) {
			sc, err := buildtestutils.GetContainer(taskDefinition, "scenv")
			require.NoError(t, err)

			scEnv := buildtestutils.KVListToMap(sc.Environment, buildtestutils.KVListToMap_KVP)
			require.Contains(t, scEnv, "OTHER_EXTRA_VAR")
			require.Contains(t, scEnv, "SC_TEST_VAR")
			require.Equal(t, "someval", scEnv["OTHER_EXTRA_VAR"])
			require.Equal(t, "overridden2", scEnv["SC_TEST_VAR"])
			require.Len(t, scEnv, 2)
		})
	})

	t.Run("port_mappings", func(t *testing.T) {
		ctx := buildtestutils.LoadProjectConfig(t, "../testdata/everything.yml")
		task := buildtestutils.GetServiceTask(ctx.Project, "svc-sidecar-ports")
		taskDefinition := genTaskDef(t, ctx, task)

		t.Run("enabled", func(t *testing.T) {
			sc1, err := buildtestutils.GetContainer(taskDefinition, "sideport")
			require.NoError(t, err)

			scPorts := buildtestutils.KVListToSlice(sc1.PortMappings, buildtestutils.KVListToSlice_PortMaps)
			require.Contains(t, scPorts, "8080/tcp")
		})

		t.Run("disabled", func(t *testing.T) {
			sc1, err := buildtestutils.GetContainer(taskDefinition, "noport")
			require.NoError(t, err)

			scPorts := buildtestutils.KVListToSlice(sc1.PortMappings, buildtestutils.KVListToSlice_PortMaps)
			require.Len(t, scPorts, 0)
		})
	})

	t.Run("depends_on", func(t *testing.T) {
		ctx := buildtestutils.LoadProjectConfig(t, "../testdata/everything.yml")
		task := buildtestutils.GetPredeployTask(ctx.Project, "pd-sc-inherit")
		taskDefinition := genTaskDef(t, ctx, task)

		sc1, err := buildtestutils.GetContainer(taskDefinition, "sc1")
		require.NoError(t, err)

		scDeps := buildtestutils.KVListToMap(sc1.DependsOn, buildtestutils.KVListToMap_Depends)
		require.Contains(t, scDeps, "scno")
		require.Equal(t, "START", scDeps["scno"])

	})
}
