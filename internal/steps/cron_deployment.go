package steps

import "ecsdeployer.com/ecsdeployer/pkg/config"

func CronDeploymentStep(resource *config.Project) *Step {

	if len(resource.CronJobs) == 0 {
		return NoopStep()
	}

	var deps []*Step
	var parallelize bool

	if resource.Settings.CronUsesEventing {
		parallelize = true

		deps = make([]*Step, len(resource.CronJobs))
		for i := range resource.CronJobs {
			deps[i] = CronjobStep(resource.CronJobs[i], true)
		}

	} else {
		parallelize = false
		deps = []*Step{
			// ensure that the group exists
			ScheduleGroupStep(resource),

			// parallelize the individual schedules
			CronSchedulesStep(resource),
		}
	}

	return NewStep(&Step{
		Label:        "CronDeployment",
		Resource:     resource,
		Dependencies: deps,
		ParallelDeps: parallelize,
	})
}
