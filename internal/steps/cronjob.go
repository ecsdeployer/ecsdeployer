package steps

import "ecsdeployer.com/ecsdeployer/pkg/config"

func CronjobStep(resource *config.CronJob) *Step {
	return NewStep(&Step{
		Label:    "Cronjob",
		ID:       resource.Name,
		Resource: resource,
		Dependencies: []*Step{
			CronTargetStep(resource),
		},
	})
}
