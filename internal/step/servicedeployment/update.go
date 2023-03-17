package servicedeployment

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	serviceBuilder "ecsdeployer.com/ecsdeployer/internal/builders/service"
	"ecsdeployer.com/ecsdeployer/internal/builders/taskdefinition"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	log "github.com/caarlos0/log"
)

func updateService(ctx *config.Context, service *config.Service) error {
	taskDef, err := taskdefinition.Register(ctx, service)
	if err != nil {
		return fmt.Errorf("failed to register task definition: %w", err)
	}

	updateServiceInput, err := serviceBuilder.BuildUpdate(ctx, service)
	if err != nil {
		return err
	}

	updateServiceInput.TaskDefinition = &taskDef.Arn

	log.WithField("name", *updateServiceInput.Service).Info("updating")
	result, err := awsclients.ECSClient().UpdateService(ctx.Context, updateServiceInput)
	if err != nil {
		return err
	}

	return waitForStable(ctx, result.Service)
}
