package steps

import (
	"errors"

	cronBuilder "ecsdeployer.com/ecsdeployer/internal/builders/cron"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
)

func CronTargetStep(resource *config.CronJob) *Step {
	return NewStep(&Step{
		Label:    "CronTarget",
		ID:       resource.Name,
		Resource: resource,
		Create:   stepCronTargetCreate,
		Dependencies: []*Step{
			TaskDefinitionStep(resource),
			CronRuleStep(resource),
		},
	})
}

func stepCronTargetCreate(ctx *config.Context, step *Step, meta *StepMetadata) (OutputFields, error) {

	taskDefinitionArn, ok := step.LookupOutput("task_definition_arn")
	if !ok {
		return nil, errors.New("Unable to find task definition arn")
	}

	// step.Logger.Info("Making cron target")

	cronJob := (step.Resource).(*config.CronJob)

	putTargetsInput, err := cronBuilder.BuildCronTarget(ctx, cronJob, taskDefinitionArn.(string))
	if err != nil {
		return nil, err
	}

	client := ctx.EventsClient()

	result, err := client.PutTargets(ctx.Context, putTargetsInput)
	if err != nil {
		return nil, err
	}

	if result.FailedEntryCount > 0 {
		for _, failEntry := range result.FailedEntries {
			step.Logger.Errorf("TargetFailure: (%s) %s: %s", aws.ToString(failEntry.TargetId), aws.ToString(failEntry.ErrorCode), aws.ToString(failEntry.ErrorMessage))
		}
		return nil, errors.New("Failed to create Cron Targets")
	}

	return nil, nil
}
