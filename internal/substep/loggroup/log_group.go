// creates/updates a log group
package loggroup

import (
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	log "github.com/caarlos0/log"
)

type Substep struct {
	groupName string
	entity    config.IsTaskStruct
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

	// figure out the group name
	tpl := tmpl.New(ctx).WithExtraFields(common.TemplateFields())
	logGroupname, err := tpl.Apply(*ctx.Project.Templates.LogGroup)
	if err != nil {
		return err
	}
	s.groupName = logGroupname

	logGroup, err := s.describeLogGroup(ctx, false)
	if err != nil {
		return err
	}

	if logGroup == nil {
		return s.createLogGroup(ctx)
	}

	return s.updateLogGroup(ctx, logGroup)
}
