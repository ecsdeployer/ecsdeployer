package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type TaskDefRegResult struct {
	Arn    string
	Family string
	Input  *ecs.RegisterTaskDefinitionInput
}

// Builds a task definition and then registers it to AWS
func Register(ctx *config.Context, entity config.IsTaskStruct) (*TaskDefRegResult, error) {

	taskInput, err := Build(ctx, entity)
	if err != nil {
		return nil, err
	}

	ecsClient := awsclients.ECSClient()

	result, err := ecsClient.RegisterTaskDefinition(ctx.Context, taskInput)
	if err != nil {
		return nil, err
	}

	taskDefArn := *result.TaskDefinition.TaskDefinitionArn

	ctx.Cache.RegisteredTaskDefArns = append(ctx.Cache.RegisteredTaskDefArns, taskDefArn)

	return &TaskDefRegResult{
		Arn:    taskDefArn,
		Family: *result.TaskDefinition.Family,
		Input:  taskInput,
	}, nil
}
