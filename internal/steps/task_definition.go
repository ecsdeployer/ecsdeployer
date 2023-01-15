package steps

import (
	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	taskBuilder "ecsdeployer.com/ecsdeployer/internal/builders/taskdef"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
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

	taskDefInput, err := taskBuilder.Build(ctx, taskStruct)
	if err != nil {
		return nil, err
	}

	ecsClient := awsclients.ECSClient()

	result, err := ecsClient.RegisterTaskDefinition(ctx.Context, taskDefInput)
	if err != nil {
		return nil, err
	}

	taskDefArn := aws.ToString(result.TaskDefinition.TaskDefinitionArn)

	step.Logger.WithField("taskDef", taskDefArn).Info("Task definition registered")

	fields := OutputFields{
		"task_definition_arn": taskDefArn,
	}

	return fields, nil
}
