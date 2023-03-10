package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func (b *Builder) applyContainerLoggingCustom(cdef *ecsTypes.ContainerDefinition, thing hasContainerAttrs) error {

	logConfig := b.project.Logging.Custom
	if logConfig.IsDisabled() {
		return nil
	}

	taskLogConfig := &config.TaskLoggingConfig{
		Driver:  util.Ptr(logConfig.Driver),
		Options: logConfig.Options.Filter(),
	}

	return b.buildContainerLogging(cdef, taskLogConfig)
}
