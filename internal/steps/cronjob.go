package steps

import (
	"errors"
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	cronBuilder "ecsdeployer.com/ecsdeployer/internal/builders/cron"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	schedulerTypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
)

const (
	cronjobStepAttrScheduleName = "CronScheduleName"
	cronjobStepAttrGroupName    = "CronGroupName"
)

func CronjobStep(resource *config.CronJob) *Step {
	return NewStep(&Step{
		Label:    "Cronjob",
		ID:       resource.Name,
		Resource: resource,
		Create:   stepCronjobCreate,
		Read:     stepCronjobRead,
		Update:   stepCronjobUpdate,
		PreApply: stepCronjobPreApply,
		Dependencies: []*Step{
			TaskDefinitionStep(resource),
		},
	})
}

func stepCronjobPreApply(ctx *config.Context, step *Step, meta *StepMetadata) error {

	common, err := config.ExtractCommonTaskAttrs(step.Resource)
	if err != nil {
		return err
	}

	tplFields, err := helpers.GetDefaultTaskTemplateFields(ctx, common)
	if err != nil {
		return err
	}

	scheduleGroupName, err := tmpl.New(ctx).Apply(*ctx.Project.Templates.ScheduleGroupName)
	if err != nil {
		return err
	}

	scheduleName, err := tmpl.New(ctx).WithExtraFields(tplFields).Apply(*ctx.Project.Templates.ScheduleName)
	if err != nil {
		return err
	}

	step.SetAttr(cronjobStepAttrScheduleName, scheduleName)
	step.SetAttr(cronjobStepAttrGroupName, scheduleGroupName)

	return nil
}

func stepCronjob_BuildCreateReq(ctx *config.Context, step *Step, meta *StepMetadata) (*scheduler.CreateScheduleInput, error) {
	taskDefinitionArn, ok := step.LookupOutput("task_definition_arn")
	if !ok {
		return nil, fmt.Errorf("%w: Unable to find task definition arn", ErrStepDependencyFailure)
	}

	cronJob := (step.Resource).(*config.CronJob)

	createScheduleInput, err := cronBuilder.BuildSchedule(ctx, cronJob, taskDefinitionArn.(string))
	if err != nil {
		return nil, err
	}

	return createScheduleInput, nil
}

func stepCronjobCreate(ctx *config.Context, step *Step, meta *StepMetadata) (OutputFields, error) {

	createInput, err := stepCronjob_BuildCreateReq(ctx, step, meta)
	if err != nil {
		return nil, err
	}

	result, err := awsclients.SchedulerClient().CreateSchedule(ctx.Context, createInput)
	if err != nil {
		return nil, err
	}
	_ = result

	return nil, nil
}

func stepCronjobUpdate(ctx *config.Context, step *Step, meta *StepMetadata) (OutputFields, error) {
	createInput, err := stepCronjob_BuildCreateReq(ctx, step, meta)
	if err != nil {
		return nil, err
	}

	updateInput := &scheduler.UpdateScheduleInput{
		FlexibleTimeWindow:         createInput.FlexibleTimeWindow,
		Name:                       createInput.Name,
		ScheduleExpression:         createInput.ScheduleExpression,
		Target:                     createInput.Target,
		ClientToken:                createInput.ClientToken,
		Description:                createInput.Description,
		EndDate:                    createInput.EndDate,
		GroupName:                  createInput.GroupName,
		KmsKeyArn:                  createInput.KmsKeyArn,
		ScheduleExpressionTimezone: createInput.ScheduleExpressionTimezone,
		StartDate:                  createInput.StartDate,
		State:                      createInput.State,
	}

	result, err := awsclients.SchedulerClient().UpdateSchedule(ctx.Context, updateInput)
	if err != nil {
		return nil, err
	}

	_ = result

	return nil, nil
}

func stepCronjobRead(ctx *config.Context, step *Step, meta *StepMetadata) (any, error) {

	scheduleGroupName := step.GetAttrMust(cronjobStepAttrGroupName).(string)
	scheduleName := step.GetAttrMust(cronjobStepAttrScheduleName).(string)

	result, err := awsclients.SchedulerClient().GetSchedule(ctx.Context, &scheduler.GetScheduleInput{
		GroupName: aws.String(scheduleGroupName),
		Name:      aws.String(scheduleName),
	})
	if err != nil {
		var rnfe *schedulerTypes.ResourceNotFoundException
		if errors.As(err, &rnfe) {
			return nil, nil
		}
		return nil, err
	}

	return result, nil
}
