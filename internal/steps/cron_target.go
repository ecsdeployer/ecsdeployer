package steps

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	cronBuilder "ecsdeployer.com/ecsdeployer/internal/builders/cron"
	"ecsdeployer.com/ecsdeployer/pkg/config"
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
		return nil, fmt.Errorf("%w: Unable to find task definition arn", ErrStepDependencyFailure)
	}

	// step.Logger.Info("Making cron target")

	cronJob := (step.Resource).(*config.CronJob)

	putTargetsInput, err := cronBuilder.BuildCronTarget(ctx, cronJob, taskDefinitionArn.(string))
	if err != nil {
		return nil, err
	}

	client := awsclients.EventsClient()

	result, err := client.PutTargets(ctx.Context, putTargetsInput)
	if err != nil {
		return nil, err
	}

	if result.FailedEntryCount > 0 {
		for _, failEntry := range result.FailedEntries {
			step.Logger.Errorf("TargetFailure: (%s) %s: %s", *failEntry.TargetId, *failEntry.ErrorCode, *failEntry.ErrorMessage)
		}
		return nil, fmt.Errorf("%w: Failed to create Cron Targets", ErrStepFailed)
	}

	return nil, nil
}
