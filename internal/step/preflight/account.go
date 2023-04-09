package preflight

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type checkAccount struct{}

func (checkAccount) String() string {
	return "aws account"
}

func (checkAccount) Skip(ctx *config.Context) bool {
	return util.IsBlank(ctx.Project.EcsDeployerOptions.AllowedAccountId)
}

func (checkAccount) Check(ctx *config.Context) error {

	accountId := ctx.Project.EcsDeployerOptions.AllowedAccountId

	if util.IsBlank(accountId) {
		return nil
	}

	if ctx.AwsAccountId() != *accountId {
		return fmt.Errorf("Account '%s' is not an allowed account. Only '%s' is allowed.", ctx.AwsAccountId(), *accountId)
	}

	return nil
}
