package taskdefinition_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
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
			router_options:
				blahwhatever: yar
				tplthing: {template: "{{.Cluster}}-stream"}
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
	})
}
