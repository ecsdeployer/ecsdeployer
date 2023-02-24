package taskdefinition

func (b *Builder) applyCleanup() error {

	// unnecessary pedantry, probably
	for i := range b.taskDef.ContainerDefinitions {
		if len(b.taskDef.ContainerDefinitions[i].Environment) == 0 {
			b.taskDef.ContainerDefinitions[i].Environment = nil
		}

		if len(b.taskDef.ContainerDefinitions[i].Secrets) == 0 {
			b.taskDef.ContainerDefinitions[i].Secrets = nil
		}

		if len(b.taskDef.ContainerDefinitions[i].DockerLabels) == 0 {
			b.taskDef.ContainerDefinitions[i].DockerLabels = nil
		}
	}

	return nil
}
