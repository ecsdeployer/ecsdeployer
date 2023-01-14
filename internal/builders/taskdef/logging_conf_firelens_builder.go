package taskdef

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func loggingConfBuilderDefaultFirelens(input *pipelineInput) (*ecsTypes.LogConfiguration, *ecsTypes.ContainerDefinition, error) {

	common := input.Common
	firelensConfig := input.Context.Project.Logging.FirelensConfig

	tpl := tmpl.New(input.Context).WithExtraFields(common.TemplateFields())

	// ctx := input.Context

	// logConfig = ctx.Project.Logging

	flConfig := &ecsTypes.FirelensConfiguration{
		Type:    ecsTypes.FirelensConfigurationType(*firelensConfig.Type),
		Options: make(map[string]string),
	}

	memory := util.Coalesce(firelensConfig.Memory)
	flContainer := &ecsTypes.ContainerDefinition{
		Name:                  firelensConfig.Name,
		Essential:             aws.Bool(true),
		Image:                 aws.String(firelensConfig.Image.Value()),
		FirelensConfiguration: flConfig,
		// MemoryReservation:     aws.Int32(int32(*memory)),
	}

	if memory != nil {
		memVal := memory.GetValueOnly()
		if memVal > 0 {
			flContainer.MemoryReservation = aws.Int32(memVal)
		}
	}

	if firelensConfig.Credentials != nil {
		flContainer.RepositoryCredentials = &ecsTypes.RepositoryCredentials{
			CredentialsParameter: firelensConfig.Credentials,
		}
	}

	if *firelensConfig.InheritEnv {
		flContainer.Environment = input.TaskDef.ContainerDefinitions[0].Environment
		flContainer.Secrets = input.TaskDef.ContainerDefinitions[0].Secrets
	}

	if firelensConfig.Options.HasSSM() {
		return nil, nil, errors.New("Cannot use SSM references in firelens options")
	}

	for lk, lv := range firelensConfig.Options.Filter() {
		val, err := lv.GetValue(tpl)
		if err != nil {
			return nil, nil, err
		}
		flConfig.Options[lk] = val
	}

	if firelensConfig.LogToAwsLogs.Enabled() {
		flContainer.LogConfiguration = &ecsTypes.LogConfiguration{
			LogDriver: ecsTypes.LogDriverAwslogs,
			Options:   make(map[string]string),
		}

		for k, v := range map[string]string{
			"awslogs-group":         firelensConfig.LogToAwsLogs.Path,
			"awslogs-region":        "{{ AwsRegion }}",
			"awslogs-stream-prefix": "firelens/{{ .Name }}",
			"awslogs-create-group":  "true",
		} {
			val, err := tpl.Apply(v)
			if err != nil {
				return nil, nil, err
			}
			flContainer.LogConfiguration.Options[k] = val
		}
	}

	conf := &ecsTypes.LogConfiguration{
		LogDriver:     ecsTypes.LogDriverAwsfirelens,
		SecretOptions: []ecsTypes.Secret{},
		Options:       make(map[string]string),
	}

	return conf, flContainer, nil

}
