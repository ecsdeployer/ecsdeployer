package service

import (
	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	serviceBuilder "ecsdeployer.com/ecsdeployer/internal/builders/service"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/webdestroya/go-log"
)

func (s *Step) updateService(ctx *config.Context) error {

	updateServiceInput, err := serviceBuilder.BuildUpdate(ctx, s.entity)
	if err != nil {
		return err
	}

	updateServiceInput.TaskDefinition = &s.taskDefArn

	log.WithField("name", s.name).Info("updating")
	result, err := awsclients.ECSClient().UpdateService(ctx.Context, updateServiceInput)
	if err != nil {
		return err
	}

	return waitForStable(ctx, result.Service)
}
