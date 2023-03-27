package pipeline

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/step/cleanup"
	"ecsdeployer.com/ecsdeployer/internal/step/console"
	"ecsdeployer.com/ecsdeployer/internal/step/crondeployment"
	"ecsdeployer.com/ecsdeployer/internal/step/predeployment"
	"ecsdeployer.com/ecsdeployer/internal/step/preflight"
	"ecsdeployer.com/ecsdeployer/internal/step/preload"
	"ecsdeployer.com/ecsdeployer/internal/step/servicedeployment"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type Stepper interface {
	fmt.Stringer
	Run(ctx *config.Context) error
}

var CleanupPipeline = []Stepper{
	cleanup.Step{},
}

var DeploymentPipeline = []Stepper{
	preflight.Step{},
	preload.Step{},
	console.Step{},
	predeployment.Step{},
	crondeployment.Step{},
	servicedeployment.Step{},
	cleanup.Step{},
}

/*
func RunPipeline(ctx *config.Context, pipeline []Stepper) error {
	for _, step := range pipeline {
		if err := skip.Maybe(
			step,
			logging.Log(
				step.String(),
				errhandler.Handle(step.Run),
			),
		)(ctx); err != nil {
			return err
		}
	}
	return nil
}
*/
