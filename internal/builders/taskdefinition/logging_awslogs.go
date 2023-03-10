package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func (b *Builder) applyContainerLoggingAwsLogs(cdef *ecsTypes.ContainerDefinition, thing hasContainerAttrs) error {

	logConfig := b.project.Logging.AwsLogConfig
	if logConfig.IsDisabled() {
		// panic("Dont disable awslogs and firelens and leave global enabled")
		// return nil, nil, nil
		return nil
	}

	logOptions := config.MergeEnvVarMaps(config.EnvVarMap{
		// "awslogs-create-group":         config.NewEnvVar(config.EnvVarTypePlain, "true"),
		"awslogs-group":         config.NewEnvVar(config.EnvVarTypeTemplated, *b.templates.LogGroup),
		"awslogs-region":        config.NewEnvVar(config.EnvVarTypeTemplated, "{{ AwsRegion }}"),
		"awslogs-stream-prefix": config.NewEnvVar(config.EnvVarTypeTemplated, *b.templates.LogStreamPrefix),
	}, logConfig.Options).Filter()

	taskLogConfig := &config.TaskLoggingConfig{
		Driver:  util.Ptr(string(ecsTypes.LogDriverAwslogs)),
		Options: logOptions,
	}

	return b.buildContainerLogging(cdef, taskLogConfig)
}
