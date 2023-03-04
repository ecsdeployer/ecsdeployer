package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func (b *Builder) applyContainerLoggingFirelens(cdef *ecsTypes.ContainerDefinition, thing hasContainerAttrs) error {

	logConfig := b.project.Logging.FirelensConfig
	if logConfig.IsDisabled() {
		return nil
	}

	taskLogConfig := &config.TaskLoggingConfig{
		Driver:  util.Ptr(string(ecsTypes.LogDriverAwsfirelens)),
		Options: logConfig.Options.Filter(),
	}

	addContainerDependency(cdef, b.loggingContainer, ecsTypes.ContainerConditionStart)

	return b.buildContainerLogging(cdef, taskLogConfig)
}
