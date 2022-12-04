package cmd

import (
	"fmt"

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

func stepRunner(options *configLoaderExtras, mode stepRunMode) error {
	ctx, cancel, err := loadProjectContext(options)

	if err != nil {
		return err
	}

	defer cancel()

	err = nil // nolint:ineffassign

	switch mode {
	case stepRunModeDeploy:
		err = stepDeploymentStepFunc(ctx.Project).Apply(ctx)
	case stepRunModeCleanup:
		err = stepCleanupOnlyStepFunc(ctx.Project).Apply(ctx)
	default:
		err = fmt.Errorf("Unknown deploy mode: %v", mode)
	}

	if err != nil {
		return err
	}

	return nil
}
