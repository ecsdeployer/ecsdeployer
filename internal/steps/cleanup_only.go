package steps

import "ecsdeployer.com/ecsdeployer/pkg/config"

func CleanupOnlyStep(project *config.Project) *Step {

	return NewStep(&Step{
		Label:    "CleanupOnly",
		Resource: project,
		Dependencies: []*Step{
			PreloadLogGroupsStep(project),
			CleanupStep(project.Settings.KeepInSync),
		},
	})
}
