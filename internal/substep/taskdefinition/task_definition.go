// registers a task definition as well as creates the log groups needed
package taskdefinition

import (
	"errors"
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/builders/taskdefinition"
	"ecsdeployer.com/ecsdeployer/internal/middleware/errhandler"
	"ecsdeployer.com/ecsdeployer/internal/middleware/skip"
	"ecsdeployer.com/ecsdeployer/internal/substep/loggroup"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type Substep struct {
	entity config.IsTaskStruct
}

func New(entity config.IsTaskStruct) *Substep {
	return &Substep{
		entity: entity,
	}
}

// will return the task definition arn
func (s *Substep) Register(ctx *config.Context) (string, error) {

	registerTaskDefInput, err := taskdefinition.Build(ctx, s.entity)
	if err != nil {
		return "", err
	}

	_ = registerTaskDefInput

	// TODO: iterate thru the task definition and find any groups?
	// try to match them against ones we know about?

	logGroupStep := loggroup.New(s.entity)
	if err := skip.Maybe(logGroupStep, errhandler.Handle(logGroupStep.Run))(ctx); err != nil {
		return "", fmt.Errorf("Failed to provision log group: %w", err)
	}

	return "X", errors.New("NOT FINISHED")
}
