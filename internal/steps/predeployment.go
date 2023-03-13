package steps

import "ecsdeployer.com/ecsdeployer/pkg/config"

func PreDeploymentStep(resource *config.Project) *Step {

	if len(resource.PreDeployTasks) == 0 {
		return NoopStep()
	}

	deps := make([]*Step, 0, len(resource.PreDeployTasks))
	wantsSharedTaskDef := false
	for i := range resource.PreDeployTasks {
		pdTask := resource.PreDeployTasks[i]
		deps = append(deps, PreDeployTaskStep(pdTask))
		if !pdTask.Disabled && pdTask.CanOverride() {
			wantsSharedTaskDef = true
		}
	}

	if wantsSharedTaskDef {
		deps = append([]*Step{pdSharedTaskDefStep(resource)}, deps...)
	}

	return NewStep(&Step{
		Label:        "PreDeployment",
		Dependencies: deps,
	})
}

func pdSharedTaskDefStep(project *config.Project) *Step {
	return TaskDefinitionStep(&config.CommonTaskAttrs{
		CommonContainerAttrs: config.CommonContainerAttrs{
			Name:    *project.Templates.SharedTaskPD,
			Command: &config.ShellCommand{"/bin/false"},
		},
	})
}
