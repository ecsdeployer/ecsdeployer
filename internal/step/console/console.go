package console

import (
	"ecsdeployer.com/ecsdeployer/internal/substep/taskdefinition"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type Step struct{}

func (Step) String() string {
	return "registering console task"
}

func (Step) Skip(ctx *config.Context) bool {
	return !ctx.Project.ConsoleTask.IsEnabled()
}

func (Step) Run(ctx *config.Context) error {

	_, err := taskdefinition.Register(ctx, ctx.Project.ConsoleTask)

	// log.WithField("arn", result.Arn).Debug("registered task definition")

	return err
}
