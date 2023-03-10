package taskdefinition

import "github.com/aws/aws-sdk-go-v2/aws"

func (b *Builder) applyRoles() error {
	if b.project.ExecutionRole != nil {
		execRoleArn, err := b.project.ExecutionRole.Arn(b.ctx)
		if err != nil {
			return err
		}
		b.taskDef.ExecutionRoleArn = aws.String(execRoleArn)
	}

	if b.project.Role != nil {
		taskRoleArn, err := b.project.Role.Arn(b.ctx)
		if err != nil {
			return err
		}
		b.taskDef.TaskRoleArn = aws.String(taskRoleArn)
	}

	return nil
}
