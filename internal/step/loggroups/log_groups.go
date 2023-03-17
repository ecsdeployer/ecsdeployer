package loggroups

import (
	"ecsdeployer.com/ecsdeployer/internal/step"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type Step struct{}

func (Step) String() string {
	return "creating log groups"
}

func (Step) Skip(ctx *config.Context) bool {
	// return len(ctx.Project.PreDeployTasks) == 0
	return false
}

func (Step) Run(ctx *config.Context) error {

	return step.Skip("NOT FINISHED")
}
