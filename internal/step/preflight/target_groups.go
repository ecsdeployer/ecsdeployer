package preflight

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type checkTargetGroups struct{}

func (checkTargetGroups) String() string {
	return "target groups"
}

func (checkTargetGroups) Check(ctx *config.Context) error {

	for _, tg := range util.DeepFindInStruct[config.TargetGroupArn](ctx.Project) {
		if _, err := tg.Name(ctx); err != nil {
			return err
		}
		if _, err := tg.Arn(ctx); err != nil {
			return err
		}
	}

	return nil
}
