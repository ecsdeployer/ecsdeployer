package servicedeployment

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	serviceBuilder "ecsdeployer.com/ecsdeployer/internal/builders/service"
	"ecsdeployer.com/ecsdeployer/internal/builders/taskdefinition"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	log "github.com/caarlos0/log"
)

func createService(ctx *config.Context, service *config.Service) error {
	taskDef, err := taskdefinition.Register(ctx, service)
	if err != nil {
		return fmt.Errorf("failed to register task definition: %w", err)
	}
	createServiceInput, err := serviceBuilder.BuildCreate(ctx, service)
	if err != nil {
		return err
	}

	createServiceInput.TaskDefinition = &taskDef.Arn

	log.WithField("name", *createServiceInput.ServiceName).Info("creating")
	result, err := awsclients.ECSClient().CreateService(ctx.Context, createServiceInput)
	if err != nil {
		return err
	}

	return waitForStable(ctx, result.Service)
}
