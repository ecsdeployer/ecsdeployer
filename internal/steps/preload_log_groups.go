package steps

import (
	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	logs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	logTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func PreloadLogGroupsStep(project *config.Project) *Step {
	return NewStep(&Step{
		Label:    "PreloadLogGroups",
		Resource: project,
		Create:   stepPreloadLogGroupsCreate,
	})
}

func stepPreloadLogGroupsCreate(ctx *config.Context, step *Step, meta *StepMetadata) (OutputFields, error) {

	if ctx.Project.Logging.IsDisabled() {
		step.Logger.Debug("AWSLogs is not desired, skipping preload")
		return nil, nil
	}

	logGroupPrefix, err := helpers.GetTemplatedPrefix(ctx, *ctx.Project.Templates.LogGroup)
	if err != nil {
		return nil, err
	}

	step.Logger.WithField("prefix", logGroupPrefix).Info("Preloading Log Group Info")

	logsClient := awsclients.LogsClient()

	request := &logs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(logGroupPrefix),
	}

	paginator := logs.NewDescribeLogGroupsPaginator(logsClient, request)

	logGroups := make([]logTypes.LogGroup, 0, ctx.Project.ApproxNumTasks())

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx.Context)
		if err != nil {
			return nil, err
		}
		logGroups = append(logGroups, output.LogGroups...)
	}

	ctx.Cache.LogGroups = logGroups

	step.Logger.WithField("prefix", logGroupPrefix).WithField("numgroups", len(logGroups)).Debug("Log preload completed")

	return nil, nil
}
