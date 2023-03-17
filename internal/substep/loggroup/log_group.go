// creates/updates a log group
package loggroup

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
	log "github.com/caarlos0/log"
)

type Substep struct {
	entity config.IsTaskStruct
}

func New(entity config.IsTaskStruct) *Substep {
	return &Substep{
		entity: entity,
	}
}

func (s *Substep) Run(ctx *config.Context) error {

	if ctx.Project.Logging.IsDisabled() || ctx.Project.Logging.AwsLogConfig.IsDisabled() {
		return nil
	}

	common := s.entity.GetCommonTaskAttrs()

	if common.LoggingConfig != nil {
		log.Trace("AwsLogs have been disabled for this task")
		// logging disabled
		return nil
	}

	logGroup, err := s.describeLogGroup(ctx)
	if err != nil {
		return err
	}

	if logGroup == nil {
		return s.createLogGroup(ctx)
	}

	return s.updateLogGroup(ctx, logGroup)
}
