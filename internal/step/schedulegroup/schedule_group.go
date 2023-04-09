package schedulegroup

import (
	"errors"
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/step"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	schedulerTypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	"github.com/webdestroya/go-log"
)

type Step struct{}

func (Step) String() string {
	return "schedule group"
}

func (Step) Run(ctx *config.Context) error {

	if ctx.Project.Settings.CronUsesEventing {
		return step.Skip("using legacy cronjob flow")
	}

	if len(ctx.Project.CronJobs) == 0 {
		return step.Skip("no cronjobs")
	}

	scheduleGroupName, err := tmpl.New(ctx).Apply(*ctx.Project.Templates.ScheduleGroupName)
	if err != nil {
		return err
	}

	ok, err := getGroup(ctx, scheduleGroupName)
	if err != nil {
		return err
	}

	if ok {
		return nil
	}

	return createGroup(ctx, scheduleGroupName)
}

func getGroup(ctx *config.Context, scheduleGroupName string) (bool, error) {
	result, err := awsclients.SchedulerClient().GetScheduleGroup(ctx.Context, &scheduler.GetScheduleGroupInput{
		Name: &scheduleGroupName,
	})
	if err != nil {
		var rnfe *schedulerTypes.ResourceNotFoundException
		if errors.As(err, &rnfe) {
			return false, nil
		}
		return false, err
	}

	if result.State != schedulerTypes.ScheduleGroupStateActive {
		return true, fmt.Errorf("schedule group is not active, but %s", result.State)
	}

	return true, nil
}

func createGroup(ctx *config.Context, scheduleGroupName string) error {

	params := &scheduler.CreateScheduleGroupInput{
		Name: &scheduleGroupName,
	}

	tagList, _, err := helpers.NameValuePair_Build_Tags(ctx, []config.NameValuePair{}, tmpl.Fields{}, schTag)
	if err != nil {
		return err
	}

	params.Tags = tagList

	if _, err := awsclients.SchedulerClient().CreateScheduleGroup(ctx.Context, params); err != nil {
		return err
	}

	log.WithField("name", scheduleGroupName).Info("created schedule group")

	return nil
}

func schTag(s1, s2 string) schedulerTypes.Tag {
	return schedulerTypes.Tag{
		Key:   &s1,
		Value: &s2,
	}
}
