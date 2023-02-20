package taskdef

import (
	"ecsdeployer.com/ecsdeployer/internal/builders/pipeline"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type TaskRolesBuilder struct{}

func (pc *TaskRolesBuilder) Apply(obj *pipeline.PipeItem[ecs.RegisterTaskDefinitionInput]) error {
	project := obj.Context.Project

	if project.ExecutionRole != nil {
		execRoleArn, err := project.ExecutionRole.Arn(obj.Context)
		if err != nil {
			return err
		}
		obj.Data.ExecutionRoleArn = aws.String(execRoleArn)
	}

	if project.Role != nil {
		taskRoleArn, err := project.Role.Arn(obj.Context)
		if err != nil {
			return err
		}
		obj.Data.TaskRoleArn = aws.String(taskRoleArn)
	}

	return nil
}
