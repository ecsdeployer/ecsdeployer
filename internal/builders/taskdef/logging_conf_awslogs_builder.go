package taskdef

import (
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func loggingConfBuilderDefaultAwslogs(input *pipelineInput) (*ecsTypes.LogConfiguration, *ecsTypes.ContainerDefinition, error) {
	ctx := input.Context
	logConfig := ctx.Project.Logging.AwsLogConfig

	if logConfig.IsDisabled() {
		panic("Dont disable awslogs and firelens and leave global enabled")
		// return nil, nil, nil
	}

	common := input.Common
	templates := ctx.Project.Templates

	tpl := tmpl.New(ctx).WithExtraFields(common.TemplateFields())

	conf := &ecsTypes.LogConfiguration{
		LogDriver:     ecsTypes.LogDriverAwslogs,
		SecretOptions: []ecsTypes.Secret{},
		Options:       make(map[string]string),
	}

	logOptions := config.MergeEnvVarMaps(map[string]config.EnvVar{
		// "awslogs-create-group":  {Value: aws.String("true")},
		"awslogs-group":         {ValueTemplate: templates.LogGroup},
		"awslogs-region":        {ValueTemplate: aws.String("{{ AwsRegion }}")},
		"awslogs-stream-prefix": {ValueTemplate: templates.LogStreamPrefix},
	}, logConfig.Options).Filter()

	for lk, lv := range logOptions {
		if lv.IsSSM() {
			conf.SecretOptions = append(conf.SecretOptions, ecsTypes.Secret{
				Name:      aws.String(lk),
				ValueFrom: lv.ValueSSM,
			})
			continue
		}

		val, err := lv.GetValue(tpl)
		if err != nil {
			return nil, nil, err
		}
		conf.Options[lk] = val

	}

	return conf, nil, nil
}
