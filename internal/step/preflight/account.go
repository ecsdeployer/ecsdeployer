package preflight

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type checkAccount struct{}

func (checkAccount) String() string {
	return "aws account"
}

func (checkAccount) Check(ctx *config.Context) error {

	if !ctx.Project.EcsDeployerOptions.IsAllowedAccountId(ctx.AwsAccountId()) {
		return fmt.Errorf("Account '%s' is not an allowed account. Only '%s' is allowed.", ctx.AwsAccountId(), *ctx.Project.EcsDeployerOptions.AllowedAccountId)
	}

	return nil
}
