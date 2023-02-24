package taskdefinition_test

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/builders/taskdefinition"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/jmespath/go-jmespath"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

/*
TESTS TO DO:

- service
- console
- cronjob
- predeploy task

- task with firelens
- task with awslogs
- task with SSM vars
- task without ssm vars
- task with/without deployment vars
*/

func TestTaskDefinitionBuilder(t *testing.T) {
	testutil.StartMocker(t, &awsmocker.MockerOptions{
		Mocks: []*awsmocker.MockedEndpoint{
			{
				Request: &awsmocker.MockedRequest{
					Service: "ecs",
					Action:  "RegisterTaskDefinition",
				},
				Response: &awsmocker.MockedResponse{
					Body: func(rr *awsmocker.ReceivedRequest) string {

						prettyJSON, _ := util.JsonifyPretty(rr.JsonPayload)
						t.Log("JSON PAYLOAD:", prettyJSON)

						taskName, _ := jmespath.Search("family", rr.JsonPayload)

						payload, _ := util.Jsonify(map[string]interface{}{
							"taskDefinition": map[string]interface{}{
								"taskDefinitionArn": fmt.Sprintf("arn:aws:ecs:%s:%s:task-definition/%s:999", rr.Region, awsmocker.DefaultAccountId, taskName.(string)),
							},
						})

						return payload
					},
				},
			},
		},
	})

	t.Run("load balanced service", func(t *testing.T) {
		ctx, err := config.NewFromYAML("testdata/dummy.yml")
		require.NoError(t, err)

		ctx.Cache.SSMSecrets = map[string]config.EnvVar{
			"SSM_VAR_1": config.NewEnvVar(config.EnvVarTypeSSM, "/fake/path/secret1"),
			"SSM_VAR_2": config.NewEnvVar(config.EnvVarTypeSSM, "/fake/path/secret2"),
			"SSM_VAR_3": config.NewEnvVar(config.EnvVarTypeSSM, "/fake/path/secret3"),
			"SSM_VAR_4": config.NewEnvVar(config.EnvVarTypeSSM, "/fake/path/secret4"),
		}
		lbService := ctx.Project.Services[0]

		builder, err := taskdefinition.NewBuilder(ctx, lbService)
		require.NoError(t, err)

		taskDefinition, err := builder.Build()
		require.NoError(t, err)

		// jsonNew, err := util.JsonifyPretty(taskDefinition)
		// require.NoError(t, err)
		// t.Log("JSON: ", jsonNew)

		_, err = awsclients.ECSClient().RegisterTaskDefinition(ctx.Context, taskDefinition)
		require.NoError(t, err)

		// taskDefinitionOld, err := Build(ctx, lbService)

		require.NotNil(t, taskDefinition.Family, "Family")
		require.EqualValues(t, "dummy-svc1", *taskDefinition.Family)

		require.NotNil(t, taskDefinition.TaskRoleArn, "TaskRoleArn")
		require.Equal(t, "arn:aws:iam::555555555555:role/faketask", *taskDefinition.TaskRoleArn)

		require.NotNil(t, taskDefinition.ExecutionRoleArn, "ExecutionRoleArn")
		require.Equal(t, "arn:aws:iam::555555555555:role/fakeexec", *taskDefinition.ExecutionRoleArn)

		require.Equal(t, ecsTypes.NetworkModeAwsvpc, taskDefinition.NetworkMode)
		require.Contains(t, taskDefinition.RequiresCompatibilities, ecsTypes.CompatibilityFargate)

		require.NotNil(t, taskDefinition.Cpu, "Cpu")
		require.Equal(t, "1024", *taskDefinition.Cpu)

		require.NotNil(t, taskDefinition.Memory, "Memory")
		require.Equal(t, "2048", *taskDefinition.Memory)

		require.GreaterOrEqual(t, len(taskDefinition.ContainerDefinitions), 1, "number of containers")
		primaryCont := taskDefinition.ContainerDefinitions[0]
		require.NotNil(t, primaryCont, "primaryCont")
		require.Equal(t, []string{"bundle", "exec", "puma", "-C", "config/puma.rb"}, primaryCont.Command)

		require.Len(t, primaryCont.PortMappings, 1)
		require.NotNil(t, primaryCont.PortMappings[0].HostPort, "HostPort")
		require.EqualValues(t, 1234, *primaryCont.PortMappings[0].HostPort)
		require.EqualValues(t, "tcp", primaryCont.PortMappings[0].Protocol)

		require.NotNil(t, primaryCont.Image, "ImageURI")
		require.Equal(t, "fake:latest", *primaryCont.Image)

		require.EqualValues(t, 0, primaryCont.Cpu)
		require.Nil(t, primaryCont.Memory)
		require.Nil(t, primaryCont.MemoryReservation)

		// jsonOld, err := util.Jsonify(taskDefinitionOld)
		// require.NoError(t, err)

		// require.JSONEq(t, jsonOld, jsonNew)

	})

	t.Run("service with firelens", func(t *testing.T) {
		ctx, err := config.NewFromYAML("testdata/firelens.yml")
		require.NoError(t, err)

		ctx.Cache.SSMSecrets = map[string]config.EnvVar{
			"SSM_VAR_1": config.NewEnvVar(config.EnvVarTypeSSM, "/fake/path/secret1"),
			"SSM_VAR_2": config.NewEnvVar(config.EnvVarTypeSSM, "/fake/path/secret2"),
			"SSM_VAR_3": config.NewEnvVar(config.EnvVarTypeSSM, "/fake/path/secret3"),
			"SSM_VAR_4": config.NewEnvVar(config.EnvVarTypeSSM, "/fake/path/secret4"),
		}
		lbService := ctx.Project.Services[0]

		builder, err := taskdefinition.NewBuilder(ctx, lbService)
		require.NoError(t, err)

		taskDefinition, err := builder.Build()
		require.NoError(t, err)

		_, err = awsclients.ECSClient().RegisterTaskDefinition(ctx.Context, taskDefinition)
		require.NoError(t, err)
	})

	t.Run("console with firelens", func(t *testing.T) {
		ctx, err := config.NewFromYAML("testdata/firelens.yml")
		require.NoError(t, err)

		ctx.Cache.SSMSecrets = map[string]config.EnvVar{
			"SSM_VAR_1": config.NewEnvVar(config.EnvVarTypeSSM, "/fake/path/secret1"),
			"SSM_VAR_2": config.NewEnvVar(config.EnvVarTypeSSM, "/fake/path/secret2"),
		}
		task := ctx.Project.ConsoleTask

		builder, err := taskdefinition.NewBuilder(ctx, task)
		require.NoError(t, err)

		taskDefinition, err := builder.Build()
		require.NoError(t, err)

		_, err = awsclients.ECSClient().RegisterTaskDefinition(ctx.Context, taskDefinition)
		require.NoError(t, err)
	})

	t.Run("console with awslogs", func(t *testing.T) {
		ctx, err := config.NewFromYAML("testdata/awslogs.yml")
		require.NoError(t, err)

		ctx.Cache.SSMSecrets = map[string]config.EnvVar{
			"SSM_VAR_1": config.NewEnvVar(config.EnvVarTypeSSM, "/fake/path/secret1"),
			"SSM_VAR_2": config.NewEnvVar(config.EnvVarTypeSSM, "/fake/path/secret2"),
		}
		task := ctx.Project.ConsoleTask

		builder, err := taskdefinition.NewBuilder(ctx, task)
		require.NoError(t, err)

		taskDefinition, err := builder.Build()
		require.NoError(t, err)

		_, err = awsclients.ECSClient().RegisterTaskDefinition(ctx.Context, taskDefinition)
		require.NoError(t, err)
	})
}
