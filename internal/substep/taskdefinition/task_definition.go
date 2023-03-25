// registers a task definition as well as creates the log groups needed
package taskdefinition

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/builders/taskdefinition"
	"ecsdeployer.com/ecsdeployer/internal/middleware/errhandler"
	"ecsdeployer.com/ecsdeployer/internal/middleware/skip"
	"ecsdeployer.com/ecsdeployer/internal/substep/loggroup"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	tieredlog "github.com/caarlos0/log"
)

type Substep struct {
	entity config.IsTaskStruct
}

func New(entity config.IsTaskStruct) *Substep {
	return &Substep{
		entity: entity,
	}
}

func Register(ctx *config.Context, entity config.IsTaskStruct) (string, error) {
	return New(entity).Register(ctx)
}

// will return the task definition arn
func (s *Substep) Register(ctx *config.Context) (string, error) {

	registerTaskDefInput, err := taskdefinition.Build(ctx, s.entity)
	if err != nil {
		return "", err
	}
	// _ = registerTaskDefInput

	// for _, container := range registerTaskDefInput.ContainerDefinitions {
	// 	if container.LogConfiguration != nil && container.LogConfiguration.LogDriver
	// }

	// TODO: iterate thru the task definition and find any groups?
	// try to match them against ones we know about?

	logGroupStep := loggroup.New(s.entity)
	if err = skip.Maybe(logGroupStep, errhandler.Handle(logGroupStep.Run))(ctx); err != nil {
		return "", fmt.Errorf("Failed to provision log group: %w", err)
	}

	result, err := awsclients.ECSClient().RegisterTaskDefinition(ctx.Context, registerTaskDefInput)
	if err != nil {
		return "", err
	}

	taskDefArn := *result.TaskDefinition.TaskDefinitionArn

	ctx.Cache.AddTaskDefinition(taskDefArn)

	tieredlog.WithField("arn", taskDefArn).Debug("registered task definition")
	return taskDefArn, nil
}
