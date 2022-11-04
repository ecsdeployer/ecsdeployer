package steps

import "ecsdeployer.com/ecsdeployer/pkg/config"

// Primary deployment step when user requests 'deploy'
func DeploymentStep(project *config.Project) *Step {

	return NewStep(&Step{
		Label:    "Deployment",
		Resource: project,
		Dependencies: []*Step{
			PreflightStep(project),
			PreloadStep(project),
			ConsoleTaskStep(project.ConsoleTask),
			PreDeploymentStep(project),
			CronDeploymentStep(project),
			ServiceDeploymentStep(project),
			DeregisterTaskDefinitionsStep(project),
			CleanupStep(project.Settings.KeepInSync),
		},
	})
}
