package preflight

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type checkRoles struct{}

func (checkRoles) String() string {
	return "roles"
}

func (checkRoles) Check(ctx *config.Context) error {

	for _, role := range util.DeepFindInStruct[config.RoleArn](ctx.Project) {
		if _, err := role.Name(ctx); err != nil {
			return err
		}
		if _, err := role.Arn(ctx); err != nil {
			return err
		}
	}

	return nil
}
