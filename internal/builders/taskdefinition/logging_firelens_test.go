package taskdefinition_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestLoggingFirelens(t *testing.T) {
	testutil.StartMocker(t, &awsmocker.MockerOptions{
		Mocks: []*awsmocker.MockedEndpoint{
			// mock_ECS_RegisterTaskDefinition_Dump(t),
			testutil.Mock_ECS_RegisterTaskDefinition_Generic(),
		},
	})

	t.Run("service with firelens", func(t *testing.T) {
		ctx := loadProjectConfig(t, "firelens.yml")
		lbService := ctx.Project.Services[0]

		taskDefinition := genTaskDef(t, ctx, lbService)
		require.NotNil(t, taskDefinition)
	})

	t.Run("console with firelens", func(t *testing.T) {
		ctx := loadProjectConfig(t, "firelens.yml", optSetNumSSMVars(4))

		task := ctx.Project.ConsoleTask

		taskDefinition := genTaskDef(t, ctx, task)
		require.NotNil(t, taskDefinition)
	})

	t.Run("everything", func(t *testing.T) {
		ctx := loadProjectConfig(t, "firelens.yml", optSetNumSSMVars(4))

		loggingYaml := `
		firelens:
			container_name: fartlog
			type: fluentbit
			credentials: /path/container/creds
			inherit_env: true
			memory: 128
			image: "custom-container/{{.Project}}:latest"
			router_options:
				blahwhatever: yar
				tplthing: {template: "{{.Cluster}}-thing"}
			options:
				region: "us-east-1"
				delivery_stream: {template: "{{.Cluster}}-stream"}
				log-driver-buffer-limit: "2097152"
				somethingsomething:
					ssm: /path/thing
		`

		loggingConf, err := yaml.ParseYAMLString[config.LoggingConfig](testutil.CleanTestYaml(loggingYaml))
		require.NoError(t, err)
		ctx.Project.Logging = loggingConf

		task := ctx.Project.ConsoleTask

		taskDefinition := genTaskDef(t, ctx, task)
		require.NotNil(t, taskDefinition)

		flContainer, err := getContainer(taskDefinition, "fartlog")
		require.NoError(t, err)
		require.NotNil(t, flContainer)

		flconfig := flContainer.FirelensConfiguration
		require.NotNil(t, flconfig)
		// firelens type
		require.Equal(t, ecsTypes.FirelensConfigurationTypeFluentbit, flconfig.Type)
		// router options
		require.Equal(t, "fakecluster-thing", flconfig.Options["tplthing"])
		require.Equal(t, "yar", flconfig.Options["blahwhatever"])

		require.NotNil(t, flContainer.RepositoryCredentials)
		require.Equal(t, "/path/container/creds", *flContainer.RepositoryCredentials.CredentialsParameter)
		require.Equal(t, "custom-container/dummy:latest", *flContainer.Image)
		require.EqualValues(t, 128, *flContainer.MemoryReservation)

		primary, err := getContainer(taskDefinition, "console")
		require.NoError(t, err)
		require.NotNil(t, primary)

		// FOR PRIMARY:
		pLogConf := primary.LogConfiguration
		require.NotNil(t, pLogConf)
		// log config driver
		require.Equal(t, ecsTypes.LogDriverAwsfirelens, pLogConf.LogDriver)
		// options
		require.Equal(t, "us-east-1", pLogConf.Options["region"])
		require.Equal(t, "fakecluster-stream", pLogConf.Options["delivery_stream"])
		require.Equal(t, "2097152", pLogConf.Options["log-driver-buffer-limit"])
		// secret options
		pLogSecrets := kvListToMap(pLogConf.SecretOptions, kvListToMap_Secret)
		require.Equal(t, "/path/thing", pLogSecrets["somethingsomething"])

		// check container dependency
		require.Condition(t, func() bool {
			if primary.DependsOn == nil || len(primary.DependsOn) == 0 {
				return false
			}
			for _, dep := range primary.DependsOn {
				if *dep.ContainerName == *flContainer.Name {
					return dep.Condition == ecsTypes.ContainerConditionStart
				}
			}
			return false
		}, "Primary container dependency")

		// hacky check
		require.Len(t, flContainer.Secrets, len(primary.Secrets))
		require.Len(t, flContainer.Environment, len(primary.Environment))

	})
}
