package steps

import "ecsdeployer.com/ecsdeployer/pkg/config"

func ConsoleTaskStep(resource *config.ConsoleTask) *Step {

	if !resource.IsEnabled() {
		return NoopStep()
	}

	return NewStep(&Step{
		Label:    "ConsoleTask",
		Resource: resource,
		Dependencies: []*Step{
			TaskDefinitionStep(resource),
		},
	})
}
