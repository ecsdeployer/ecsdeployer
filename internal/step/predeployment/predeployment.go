package predeployment

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type Step struct{}

func (Step) String() string {
	return "running predeploy tasks"
}

func (Step) Skip(ctx *config.Context) bool {
	return len(ctx.Project.PreDeployTasks) == 0
}

func (Step) Run(ctx *config.Context) error {

	return nil
}
