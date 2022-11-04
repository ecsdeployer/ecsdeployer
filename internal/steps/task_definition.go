package steps

import (
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

	ecsClient := ctx.ECSClient()

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

/*
func stepTaskDefinitionPreApply(ctx *config.Context, step *Step, meta *StepMetadata) error {

	if ctx.Project.Logging.IsDisabled() {
		// they don't want any logging
		return nil
	}

	taskStruct := (step.Resource).(config.IsTaskStruct)

	logGroupPrefix, err := helpers.GetTemplatedPrefix(ctx, ctx.Project.NameTemplates.LogGroup)
	if err != nil {
		return err
	}

	taskDefInput, err := builders.BuildTaskDefinition(ctx, taskStruct)
	if err != nil {
		return err
	}

	for _, containerDef := range taskDefInput.ContainerDefinitions {
		if containerDef.LogConfiguration != nil && containerDef.LogConfiguration.LogDriver == ecstypes.LogDriverAwslogs {
			desiredLogGroup, ok := containerDef.LogConfiguration.Options["awslogs-group"]
			if ok && strings.HasPrefix(desiredLogGroup, logGroupPrefix) {
				// it wants a log group, and it has a prefix we desire. We need to create it
				step.Dependencies = append(step.Dependencies, Lo)
			}
		}
	}

	return nil
}
*/
