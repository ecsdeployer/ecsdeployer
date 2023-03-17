package pipeline

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/step/cleanup"
	"ecsdeployer.com/ecsdeployer/internal/step/console"
	"ecsdeployer.com/ecsdeployer/internal/step/crondeployment"
	"ecsdeployer.com/ecsdeployer/internal/step/loggroups"
	"ecsdeployer.com/ecsdeployer/internal/step/predeployment"
	"ecsdeployer.com/ecsdeployer/internal/step/preflight"
	"ecsdeployer.com/ecsdeployer/internal/step/preload"
	"ecsdeployer.com/ecsdeployer/internal/step/servicedeployment"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	// "ecsdeployer.com/ecsdeployer/internal/step/taskdefinitions"
)

type Stepper interface {
	fmt.Stringer

	// Run the pipe
	Run(ctx *config.Context) error
}

var CleanupPipeline = []Stepper{
	cleanup.Step{},
}

var DeploymentPipeline = []Stepper{
	preflight.Step{},
	preload.Step{},
	loggroups.Step{},
	console.Step{},
	// taskdefinitions.Step{},
	predeployment.Step{},
	crondeployment.Step{},
	servicedeployment.Step{},
	cleanup.Step{},
}
