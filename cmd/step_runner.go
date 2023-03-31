package cmd

import (
	"ecsdeployer.com/ecsdeployer/internal/steps"
)

type stepRunMode int

const (
	stepRunModeDeploy stepRunMode = iota
	stepRunModeCleanup
)

var (
	stepDeploymentStepFunc  = steps.DeploymentStep
	stepCleanupOnlyStepFunc = steps.CleanupOnlyStep
)
