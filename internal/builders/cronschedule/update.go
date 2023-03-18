package cronschedule

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
)

func BuildUpdate(ctx *config.Context, resource *config.CronJob, taskDefArn string) (*scheduler.UpdateScheduleInput, error) {
	createInput, err := BuildCreate(ctx, resource, taskDefArn)
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

	return updateInput, nil
}
