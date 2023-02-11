package steps

import "ecsdeployer.com/ecsdeployer/pkg/config"

func CronDeploymentStep(resource *config.Project) *Step {

	if len(resource.CronJobs) == 0 {
		return NoopStep()
	}

	// Uses the old EventBridge system of rule/targets
	if resource.Settings.CronUsesEventing {
		deps := make([]*Step, len(resource.CronJobs))
		for i := range resource.CronJobs {
			deps[i] = CronjobStep(resource.CronJobs[i], true)
		}

		return NewStep(&Step{
			Label:        "CronDeployment",
			Resource:     resource,
			Dependencies: deps,
			ParallelDeps: true,
		})
	}

	return NewStep(&Step{
		Label:    "CronDeployment",
		Resource: resource,
		Dependencies: []*Step{
			// ensure that the group exists
			ScheduleGroupStep(resource),

			// parallelize the individual schedules
			CronSchedulesStep(resource),
		},
	})
}
