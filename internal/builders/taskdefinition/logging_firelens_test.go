package taskdefinition_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/testutil/buildtestutils"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
)

func TestLoggingFirelens(t *testing.T) {

	buildtestutils.StartMocker(t)

	t.Run("service with firelens", func(t *testing.T) {
		ctx := buildtestutils.LoadProjectConfig(t, "../testdata/firelens.yml")
		lbService := ctx.Project.Services[0]

		taskDefinition := genTaskDef(t, ctx, lbService)
		require.NotNil(t, taskDefinition)

		container, _ := buildtestutils.GetContainer(taskDefinition, lbService.Name)
		require.NotNil(t, container.LogConfiguration)
		require.Equal(t, ecsTypes.LogDriverAwsfirelens, container.LogConfiguration.LogDriver)
	})

	t.Run("console with firelens", func(t *testing.T) {
		ctx := buildtestutils.LoadProjectConfig(t, "../testdata/firelens.yml", buildtestutils.OptSetNumSSMVars(4))

		task := ctx.Project.ConsoleTask

		taskDefinition := genTaskDef(t, ctx, task)
		require.NotNil(t, taskDefinition)

		container, _ := buildtestutils.GetContainer(taskDefinition, "console")
		require.NotNil(t, container.LogConfiguration)
		require.Equal(t, ecsTypes.LogDriverAwsfirelens, container.LogConfiguration.LogDriver)
	})

	t.Run("everything", func(t *testing.T) {
		ctx := buildtestutils.LoadProjectConfig(t, "../testdata/firelens.yml", buildtestutils.OptSetNumSSMVars(4))

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
				"@type": "cwlogs"
				somethingsomething:
					ssm: /path/thing
		`

		loggingConf, err := yaml.ParseYAMLString[config.LoggingConfig](testutil.CleanTestYaml(loggingYaml))
		require.NoError(t, err)
		ctx.Project.Logging = loggingConf

		task := ctx.Project.ConsoleTask

		taskDefinition := genTaskDef(t, ctx, task)
		require.NotNil(t, taskDefinition)

		flContainer, err := buildtestutils.GetContainer(taskDefinition, "fartlog")
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

		primary, err := buildtestutils.GetContainer(taskDefinition, "console")
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
		pLogSecrets := buildtestutils.KVListToMap(pLogConf.SecretOptions, buildtestutils.KVListToMap_Secret)
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

	t.Run("firelens with awslogs", func(t *testing.T) {
		ctx := buildtestutils.LoadProjectConfig(t, "../testdata/firelens.yml", buildtestutils.OptSetNumSSMVars(4))

		loggingYaml := `
		firelens:
			log_to_awslogs: /something/something
		`

		loggingConf, err := yaml.ParseYAMLString[config.LoggingConfig](testutil.CleanTestYaml(loggingYaml))
		require.NoError(t, err)
		ctx.Project.Logging = loggingConf

		task := ctx.Project.ConsoleTask

		taskDefinition := genTaskDef(t, ctx, task)
		require.NotNil(t, taskDefinition)

		flContainer, err := buildtestutils.GetContainer(taskDefinition, "log_router")
		require.NoError(t, err)
		require.NotNil(t, flContainer)

		require.NotNil(t, flContainer.LogConfiguration)
		require.Equal(t, ecsTypes.LogDriverAwslogs, flContainer.LogConfiguration.LogDriver)
		require.Equal(t, "/something/something", flContainer.LogConfiguration.Options["awslogs-group"])
		require.Equal(t, "firelens-console", flContainer.LogConfiguration.Options["awslogs-stream-prefix"])
	})

	t.Run("firelens with awslogs stream hack", func(t *testing.T) {
		ctx := buildtestutils.LoadProjectConfig(t, "../testdata/firelens.yml", buildtestutils.OptSetNumSSMVars(4))

		loggingYaml := `
		firelens:
			log_to_awslogs: "/path/{{.Project}}/flogs/etc:hacky-{{.Name}}"
		`

		loggingConf, err := yaml.ParseYAMLString[config.LoggingConfig](testutil.CleanTestYaml(loggingYaml))
		require.NoError(t, err)
		ctx.Project.Logging = loggingConf

		task := ctx.Project.ConsoleTask

		taskDefinition := genTaskDef(t, ctx, task)
		require.NotNil(t, taskDefinition)

		flContainer, err := buildtestutils.GetContainer(taskDefinition, "log_router")
		require.NoError(t, err)
		require.NotNil(t, flContainer)

		require.NotNil(t, flContainer.LogConfiguration)
		require.Equal(t, ecsTypes.LogDriverAwslogs, flContainer.LogConfiguration.LogDriver)
		require.Equal(t, "/path/dummy/flogs/etc", flContainer.LogConfiguration.Options["awslogs-group"])
		require.Equal(t, "hacky-console", flContainer.LogConfiguration.Options["awslogs-stream-prefix"])

	})

	t.Run("firelens without awslogs", func(t *testing.T) {
		ctx := buildtestutils.LoadProjectConfig(t, "../testdata/firelens.yml", buildtestutils.OptSetNumSSMVars(4))

		loggingYaml := `
		firelens:
			type: fluentbit
		`

		loggingConf, err := yaml.ParseYAMLString[config.LoggingConfig](testutil.CleanTestYaml(loggingYaml))
		require.NoError(t, err)
		ctx.Project.Logging = loggingConf

		task := ctx.Project.ConsoleTask

		taskDefinition := genTaskDef(t, ctx, task)
		require.NotNil(t, taskDefinition)

		flContainer, err := buildtestutils.GetContainer(taskDefinition, "log_router")
		require.NoError(t, err)
		require.NotNil(t, flContainer)

		require.Nil(t, flContainer.LogConfiguration)

	})
}
