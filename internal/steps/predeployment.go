package steps

import "ecsdeployer.com/ecsdeployer/pkg/config"

func PreDeploymentStep(resource *config.Project) *Step {

	if len(resource.PreDeployTasks) == 0 {
		return NoopStep()
	}

	deps := make([]*Step, 0, len(resource.PreDeployTasks))
	for i := range resource.PreDeployTasks {
		deps = append(deps, PreDeployTaskStep(resource.PreDeployTasks[i]))
	}

	return NewStep(&Step{
		Label:        "PreDeployment",
		Dependencies: deps,
	})
}
