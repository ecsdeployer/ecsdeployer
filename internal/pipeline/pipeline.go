package pipeline

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/step/cleanup"
	"ecsdeployer.com/ecsdeployer/internal/step/console"
	"ecsdeployer.com/ecsdeployer/internal/step/crondeployment"
	"ecsdeployer.com/ecsdeployer/internal/step/deregistertaskdefinitions"
	"ecsdeployer.com/ecsdeployer/internal/step/predeployment"
	"ecsdeployer.com/ecsdeployer/internal/step/preflight"
	"ecsdeployer.com/ecsdeployer/internal/step/preload"
	"ecsdeployer.com/ecsdeployer/internal/step/schedulegroup"
	"ecsdeployer.com/ecsdeployer/internal/step/servicedeployment"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type Stepper interface {
	fmt.Stringer
	Run(ctx *config.Context) error
}

var CleanupPipeline = []Stepper{
	preflight.Step{},
	cleanup.Step{},
}

var DeploymentPipeline = []Stepper{
	preflight.Step{},
	preload.Step{},
	console.Step{},
	predeployment.Step{},
	schedulegroup.Step{},
	crondeployment.Step{},
	servicedeployment.Step{},
	deregistertaskdefinitions.Step{},
	cleanup.Step{},
}
