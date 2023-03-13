package steps

import (
	taskBuilder "ecsdeployer.com/ecsdeployer/internal/builders/taskdefinition"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func TaskDefinitionStep(resource config.IsTaskStruct) *Step {

	common, err := config.ExtractCommonTaskAttrs(resource)
	if err != nil {
		panic(err)
	}

	return NewStep(&Step{
		Label:    "TaskDefinition",
		ID:       common.Name,
		Resource: resource,
		Create:   stepTaskDefinitionCreate,
		Dependencies: []*Step{
			LogGroupStep(common),
		},
	})
}

func stepTaskDefinitionCreate(ctx *config.Context, step *Step, meta *StepMetadata) (OutputFields, error) {

	taskStruct := (step.Resource).(config.IsTaskStruct)

	taskDefResult, err := taskBuilder.Register(ctx, taskStruct)
	if err != nil {
		return nil, err
	}

	step.Logger.WithField("taskDef", taskDefResult.Arn).Info("Task definition registered")

	fields := OutputFields{
		"task_definition_arn": taskDefResult.Arn,
	}

	return fields, nil
}
