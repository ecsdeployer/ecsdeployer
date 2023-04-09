package cleanupcronjobs

import (
	"fmt"
	"strings"
	"sync"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/step"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	scheduler "github.com/aws/aws-sdk-go-v2/service/scheduler"
	"github.com/webdestroya/go-log"
	"golang.org/x/exp/slices"
)

func runSchedulerCleanup(ctx *config.Context) error {

	scheduleGroupName, err := tmpl.New(ctx).Apply(*ctx.Project.Templates.ScheduleGroupName)
	if err != nil {
		return fmt.Errorf("failed to determine schedule group name: %w", err)
	}

	schedulePrefix, err := helpers.GetTemplatedPrefix(ctx, *ctx.Project.Templates.ScheduleName)
	if err != nil {
		return err
	}

	expectedScheduleNames := make([]string, 0, len(ctx.Project.CronJobs))
	for _, cron := range ctx.Project.CronJobs {
		scheduleName, err := tmpl.New(ctx).WithExtraFields(cron.TemplateFields()).Apply(*ctx.Project.Templates.ScheduleName)
		if err != nil {
			// log.WithError(err).Error("Unable to determine existing cron names. Cronjobs will not be synced")
			return step.Skipf("unable to determine cron names?? error: %s", err.Error())
		}
		expectedScheduleNames = append(expectedScheduleNames, scheduleName)
	}

	client := awsclients.SchedulerClient()

	request := &scheduler.ListSchedulesInput{
		GroupName: &scheduleGroupName,
	}

	paginator := scheduler.NewListSchedulesPaginator(client, request)

	// TODO: move this to a semerrgroup?
	var wg sync.WaitGroup

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx.Context)
		if err != nil {
			log.WithError(err).Warn("Failed to page schedule resources")
			return nil
		}
		for _, result := range output.Schedules {

			scheduleName := *result.Name

			if !strings.HasPrefix(scheduleName, schedulePrefix) {
				// this schedule doesnt follow our name convention, assume it was not created by us
				log.WithField("schedule", scheduleName).Trace("ignoring non-ecsdeploy schedule")
				continue
			}

			if slices.Contains(expectedScheduleNames, scheduleName) {
				// service is supposed to be there, so dont delete
				continue
			}

			wg.Add(1)

			go func(groupName, schedName string) {
				defer wg.Done()
				if err := destroySchedule(ctx, groupName, schedName); err != nil {
					log.WithError(err).WithField("schedule", schedName).Warn("unable to delete (ignoring)")
				}
			}(scheduleGroupName, scheduleName)
		}
	}

	wg.Wait()

	return nil
}

func destroySchedule(ctx *config.Context, groupName, scheduleName string) error {

	log.WithField("schedule", scheduleName).Info("found unwanted schedule")
	client := awsclients.SchedulerClient()

	_, err := client.DeleteSchedule(ctx.Context, &scheduler.DeleteScheduleInput{
		Name:      &scheduleName,
		GroupName: &groupName,
	})
	if err != nil {
		log.WithField("schedule", scheduleName).WithError(err).Warn("failed to delete schedule")
		return err
	}
	log.WithField("schedule", scheduleName).Info("deleted schedule")

	return nil
}
