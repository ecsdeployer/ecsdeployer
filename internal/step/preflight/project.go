package preflight

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type checkProject struct{}

func (checkProject) String() string {
	return "project config"
}

func (checkProject) Check(ctx *config.Context) error {
	return ctx.Project.ValidateWithContext(ctx)
}
