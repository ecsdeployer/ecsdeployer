package taskdefinition_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestSidecars(t *testing.T) {
	testutil.StartMocker(t, &awsmocker.MockerOptions{
		Mocks: []*awsmocker.MockedEndpoint{
			// Mock_ECS_RegisterTaskDefinition_Dump(t),
			testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
		},
	})

	t.Run("inherit_env", func(t *testing.T) {
		ctx := loadProjectConfig(t, "everything.yml")

		task := getPredeployTask(ctx.Project, "pd-sc-inherit")
		taskDefinition := genTaskDef(t, ctx, task)

		t.Run("enabled", func(t *testing.T) {
			sc1, err := getContainer(taskDefinition, "sc1")
			require.NoError(t, err)

			scEnv := kvListToMap(sc1.Environment, kvListToMap_KVP)
			require.Contains(t, scEnv, "SC_TEST_VAR")
			require.Equal(t, "blah", scEnv["SC_TEST_VAR"])
		})

		t.Run("disabled", func(t *testing.T) {
			sc1, err := getContainer(taskDefinition, "scno")
			require.NoError(t, err)

			scEnv := kvListToMap(sc1.Environment, kvListToMap_KVP)
			require.NotContains(t, scEnv, "SC_TEST_VAR")
		})

		t.Run("default", func(t *testing.T) {
			sc1, err := getContainer(taskDefinition, "scdef")
			require.NoError(t, err)

			scEnv := kvListToMap(sc1.Environment, kvListToMap_KVP)
			require.NotContains(t, scEnv, "SC_TEST_VAR")
		})
	})

	t.Run("port_mappings", func(t *testing.T) {
		ctx := loadProjectConfig(t, "everything.yml")
		task := getServiceTask(ctx.Project, "svc-sidecar-ports")
		taskDefinition := genTaskDef(t, ctx, task)

		t.Run("enabled", func(t *testing.T) {
			sc1, err := getContainer(taskDefinition, "sideport")
			require.NoError(t, err)

			scPorts := kvListToSlice(sc1.PortMappings, kvListToSlice_PortMaps)
			require.Contains(t, scPorts, "8080/tcp")
		})

		t.Run("disabled", func(t *testing.T) {
			sc1, err := getContainer(taskDefinition, "noport")
			require.NoError(t, err)

			scPorts := kvListToSlice(sc1.PortMappings, kvListToSlice_PortMaps)
			require.Len(t, scPorts, 0)
		})
	})

	t.Run("depends_on", func(t *testing.T) {
		ctx := loadProjectConfig(t, "everything.yml")
		task := getPredeployTask(ctx.Project, "pd-sc-inherit")
		taskDefinition := genTaskDef(t, ctx, task)

		sc1, err := getContainer(taskDefinition, "sc1")
		require.NoError(t, err)

		scDeps := kvListToMap(sc1.DependsOn, kvListToMap_Depends)
		require.Contains(t, scDeps, "scno")
		require.Equal(t, "START", scDeps["scno"])

	})
}
