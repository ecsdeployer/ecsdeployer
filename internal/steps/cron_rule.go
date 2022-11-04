package steps

import (
	cronBuilder "ecsdeployer.com/ecsdeployer/internal/builders/cron"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func CronRuleStep(resource *config.CronJob) *Step {
	return NewStep(&Step{
		Label:    "CronRule",
		ID:       resource.Name,
		Resource: resource,
		Create:   stepCronRuleCreate,
	})
}

func stepCronRuleCreate(ctx *config.Context, step *Step, meta *StepMetadata) (OutputFields, error) {

	cronjob := (step.Resource).(*config.CronJob)

	payload, err := cronBuilder.BuildCronRule(ctx, cronjob)
	if err != nil {
		return nil, err
	}

	step.Logger.WithField("rule", *payload.Name).Info("Creating CronJob Rule")
	result, err := ctx.EventsClient().PutRule(ctx, payload)
	if err != nil {
		return nil, err
	}

	fields := OutputFields{
		"rule_arn": *result.RuleArn,
	}

	return fields, nil
}
