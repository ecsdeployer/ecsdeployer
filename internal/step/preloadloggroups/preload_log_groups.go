package preloadloggroups

import (
	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
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

	logsClient := awsclients.LogsClient()

	request := &logs.DescribeLogGroupsInput{
		LogGroupNamePrefix: &logGroupPrefix,
	}

	paginator := logs.NewDescribeLogGroupsPaginator(logsClient, request)

	logGroups := make([]logTypes.LogGroup, 0, ctx.Project.ApproxNumTasks())

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx.Context)
		if err != nil {
			return err
		}
		logGroups = append(logGroups, output.LogGroups...)
	}

	ctx.Cache.LogGroups = logGroups

	return nil
}
