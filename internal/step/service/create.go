package service

import (
	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	serviceBuilder "ecsdeployer.com/ecsdeployer/internal/builders/service"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/webdestroya/go-log"
)

func (s *Step) createService(ctx *config.Context) error {

	createServiceInput, err := serviceBuilder.BuildCreate(ctx, s.entity)
	if err != nil {
		return err
	}

	createServiceInput.TaskDefinition = &s.taskDefArn

	log.WithField("name", s.name).Info("creating")
	result, err := awsclients.ECSClient().CreateService(ctx.Context, createServiceInput)
	if err != nil {
		return err
	}

	return waitForStable(ctx, result.Service)
}
