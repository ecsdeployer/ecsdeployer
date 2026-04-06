package taskdefinition

func (b *Builder) applyRoles() error {
	if b.project.ExecutionRole != nil {
		execRoleArn, err := b.project.ExecutionRole.Arn(b.ctx)
		if err != nil {
			return err
		}
		b.taskDef.ExecutionRoleArn = new(execRoleArn)
	}

	if b.project.Role != nil {
		taskRoleArn, err := b.project.Role.Arn(b.ctx)
		if err != nil {
			return err
		}
		b.taskDef.TaskRoleArn = new(taskRoleArn)
	}

	return nil
}
