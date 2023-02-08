package steps

import "ecsdeployer.com/ecsdeployer/pkg/config"

func CronDeploymentStep(resource *config.Project) *Step {

	if len(resource.CronJobs) == 0 {
		return NoopStep()
	}

	deps := []*Step{
		// ensure that the group exists
		ScheduleGroupStep(resource),

		// parallelize the individual schedules
		CronSchedulesStep(resource),
	}

	return NewStep(&Step{
		Label:        "CronDeployment",
		Resource:     resource,
		Dependencies: deps,
	})
}
