package cronjob

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	scheduleBuilder "ecsdeployer.com/ecsdeployer/internal/builders/cronschedule"
	"ecsdeployer.com/ecsdeployer/internal/substep/taskdefinition"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/scheduler"
	schedulerTypes "github.com/aws/aws-sdk-go-v2/service/scheduler/types"
	log "github.com/caarlos0/log"
)

func (s *Step) runSchedulerVariant(ctx *config.Context) error {

	scheduleGroupName, err := tmpl.New(ctx).Apply(*ctx.Project.Templates.ScheduleGroupName)
	if err != nil {
		return err
	}
	s.groupName = scheduleGroupName

	scheduleName, err := tmpl.New(ctx).WithExtraFields(s.cronjob.TemplateFields()).Apply(*ctx.Project.Templates.ScheduleName)
	if err != nil {
		return err
	}
	s.scheduleName = scheduleName

	exists, err := s.getSchedule(ctx)
	if err != nil {
		return err
	}

	taskDefArn, err := taskdefinition.Register(ctx, s.cronjob)
	if err != nil {
		return err
	}
	s.taskDefArn = taskDefArn

	if exists {
		return s.updateSchedule(ctx)
	}

	return s.createSchedule(ctx)
}

func (s *Step) getSchedule(ctx *config.Context) (bool, error) {
	_, err := awsclients.SchedulerClient().GetSchedule(ctx.Context, &scheduler.GetScheduleInput{
		GroupName: &s.groupName,
		Name:      &s.scheduleName,
	})
	if err != nil {
		var rnfe *schedulerTypes.ResourceNotFoundException
		if errors.As(err, &rnfe) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *Step) createSchedule(ctx *config.Context) error {

	input, err := scheduleBuilder.BuildCreate(ctx, s.cronjob, s.taskDefArn)
	if err != nil {
		return err
	}

	if _, err := awsclients.SchedulerClient().CreateSchedule(ctx.Context, input); err != nil {
		return err
	}
	log.WithField("name", s.cronjob.Name).Info("created schedule")

	return nil
}

func (s *Step) updateSchedule(ctx *config.Context) error {
	input, err := scheduleBuilder.BuildUpdate(ctx, s.cronjob, s.taskDefArn)
	if err != nil {
		return err
	}

	if _, err := awsclients.SchedulerClient().UpdateSchedule(ctx.Context, input); err != nil {
		return err
	}

	log.WithField("name", s.cronjob.Name).Info("updated schedule")

	return nil
}
