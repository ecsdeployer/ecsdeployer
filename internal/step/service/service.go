package service

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/substep/taskdefinition"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type Step struct {
	entity     *config.Service
	name       string
	taskDefArn string
}

func New(svc *config.Service) *Step {
	return &Step{
		entity: svc,
	}
}

func (Step) String() string {
	return "deploying services"
}

func (s *Step) Run(ctx *config.Context) error {

	serviceName, err := tmpl.New(ctx).WithExtraFields(s.entity.TemplateFields()).Apply(*ctx.Project.Templates.ServiceName)
	if err != nil {
		return err
	}
	s.name = serviceName

	taskDefArn, err := taskdefinition.Register(ctx, s.entity)
	if err != nil {
		return fmt.Errorf("failed to register task definition: %w", err)
	}
	s.taskDefArn = taskDefArn

	exists, err := s.getExisting(ctx)
	if err != nil {
		return err
	}

	if exists {
		return s.updateService(ctx)
	}

	return s.createService(ctx)
}
