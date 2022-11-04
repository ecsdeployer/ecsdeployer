package steps

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func PreloadStep(project *config.Project) *Step {
	return NewStep(&Step{
		Label:        "Preload",
		Resource:     project,
		ParallelDeps: true,

		Dependencies: []*Step{
			PreloadSecretsStep(project),
			PreloadLogGroupsStep(project),
		},
	})
}
