package steps

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	schedulerTypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
)

func ScheduleGroupStep(resource *config.Project) *Step {
	if len(resource.CronJobs) == 0 {
		return NoopStep()
	}

	return NewStep(&Step{
		Label:    "ScheduleGroup",
		Resource: resource,
		PreApply: stepScheduleGroupPreApply,
		Read:     stepScheduleGroupRead,
		Create:   stepScheduleGroupCreate,
	})
}

func stepScheduleGroupPreApply(ctx *config.Context, step *Step, meta *StepMetadata) error {
	tpl := tmpl.New(ctx)

	scheduleGroupName, err := tpl.Apply(*ctx.Project.Templates.ScheduleGroupName)
	if err != nil {
		return err
	}

	step.Attributes["scheduleGroupName"] = scheduleGroupName

	return nil
}

func stepScheduleGroupRead(ctx *config.Context, step *Step, meta *StepMetadata) (any, error) {

	scheduleGroupName := step.Attributes["scheduleGroupName"].(string)

	result, err := awsclients.SchedulerClient().GetScheduleGroup(ctx.Context, &scheduler.GetScheduleGroupInput{
		Name: aws.String(scheduleGroupName),
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

func stepScheduleGroupCreate(ctx *config.Context, step *Step, meta *StepMetadata) (OutputFields, error) {

	scheduleGroupName := step.Attributes["scheduleGroupName"].(string)

	params := &scheduler.CreateScheduleGroupInput{
		Name: aws.String(scheduleGroupName),
		Tags: []schedulerTypes.Tag{},
	}

	tagList, _, err := helpers.NameValuePair_Build_Tags(ctx, []config.NameValuePair{}, tmpl.Fields{}, func(s1, s2 string) schedulerTypes.Tag {
		return schedulerTypes.Tag{
			Key:   &s1,
			Value: &s2,
		}
	})
	if err != nil {
		return nil, err
	}

	params.Tags = tagList

	result, err := awsclients.SchedulerClient().CreateScheduleGroup(ctx.Context, params)
	if err != nil {
		return nil, err
	}

	fields := OutputFields{
		"ScheduleGroupArn": result.ScheduleGroupArn,
	}

	return fields, nil
}
