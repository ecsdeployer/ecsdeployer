package preloadloggroups

import (
	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/step"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	logs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	logTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

type Step struct{}

func (Step) String() string {
	return "preloading log groups"
}

func (Step) Skip(ctx *config.Context) bool {
	return ctx.Project.Logging.IsDisabled()
}

func (Step) Preload(ctx *config.Context) error {

	logGroupPrefix, err := helpers.GetTemplatedPrefix(ctx, *ctx.Project.Templates.LogGroup)
	if err != nil {
		return err
	}

	if util.IsBlank(&logGroupPrefix) {
		return step.Skip("Log group prefix was blank, will not attempt to cache groups")
	}

	if err := ByPrefix(ctx, logGroupPrefix); err != nil {
		return err
	}

	ctx.Cache.LogGroupsCached = true

	return nil
}

func ByPrefix(ctx *config.Context, prefix string) error {
	if ctx.Cache.LogGroups == nil {
		ctx.Cache.LogGroups = make(map[string]logTypes.LogGroup)
	}
	logsClient := awsclients.LogsClient()

	request := &logs.DescribeLogGroupsInput{
		LogGroupNamePrefix: &prefix,
	}

	paginator := logs.NewDescribeLogGroupsPaginator(logsClient, request)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx.Context)
		if err != nil {
			return err
		}
		for _, logGroup := range output.LogGroups {
			ctx.Cache.LogGroups[*logGroup.LogGroupName] = logGroup
		}
	}

	return nil
}
