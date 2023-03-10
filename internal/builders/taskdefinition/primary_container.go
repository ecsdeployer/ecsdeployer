package taskdefinition

import (
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func (b *Builder) applyPrimaryContainer() error {

	b.primaryContainer = &ecsTypes.ContainerDefinition{}

	mergedCommon := b.entity.GetCommonContainerAttrs().NewDefaultedBy(b.taskDefaults.CommonContainerAttrs)
	mergedCommon.Cpu = nil
	mergedCommon.Memory = nil

	if err := b.applyContainerDefaults(b.primaryContainer, &mergedCommon); err != nil {
		return err
	}
	if err := b.addEnvVarsToContainer(b.primaryContainer, b.baseEnvVars); err != nil {
		return err
	}

	if err := b.applyContainerLogging(b.primaryContainer, b.entity); err != nil {
		return err
	}

	// if err := b.applyContainerHealthCheck(b.primaryContainer, util.Coalesce(b.commonTask.HealthCheck, b.taskDefaults.HealthCheck)); err != nil {
	// 	return err
	// }

	return nil
}
