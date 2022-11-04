package steps

import "ecsdeployer.com/ecsdeployer/pkg/config"

func CronDeploymentStep(resource *config.Project) *Step {

	if len(resource.CronJobs) == 0 {
		return NoopStep()
	}

	deps := make([]*Step, len(resource.CronJobs))
	for i := range resource.CronJobs {
		deps[i] = CronjobStep(resource.CronJobs[i])
	}

	return NewStep(&Step{
		Label:        "CronDeployment",
		Resource:     resource,
		ParallelDeps: true, // these do not depend on each other
		Dependencies: deps,
	})
}
