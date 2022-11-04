package steps

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func ServiceDeploymentStep(project *config.Project) *Step {

	if len(project.Services) == 0 {
		return NoopStep()
	}

	depList := make([]*Step, 0, len(project.Services))
	for i := range project.Services {
		depList = append(depList, ServiceStep(project.Services[i]))
	}

	return NewStep(&Step{
		Label:        "ServiceDeployment",
		Resource:     project,
		ParallelDeps: true,
		Dependencies: depList,
	})
}
