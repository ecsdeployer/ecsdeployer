package taskdefinition

import (
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func (b *Builder) applyPrimaryContainer() error {

	b.primaryContainer = &ecsTypes.ContainerDefinition{
		// Command:                []string{},
		// Cpu:                    0,
		// DependsOn:              []ecsTypes.ContainerDependency{},
		// DockerLabels:           map[string]string{},
		// EntryPoint:             []string{},
		// Environment:            []ecsTypes.KeyValuePair{},
		// EnvironmentFiles:       []ecsTypes.EnvironmentFile{},
		// Essential:              new(bool),
		// FirelensConfiguration:  &ecsTypes.FirelensConfiguration{},
		// HealthCheck:            &ecsTypes.HealthCheck{},
		// Image:                  new(string),
		// Interactive:            new(bool),
		// LinuxParameters:        &ecsTypes.LinuxParameters{},
		// LogConfiguration: 			 &ecsTypes.LogConfiguration{},
		// Memory:                 new(int32),
		// MemoryReservation:      new(int32),
		// MountPoints:            []ecsTypes.MountPoint{},
		// Name:                   new(string),
		// PortMappings:           []ecsTypes.PortMapping{},
		// RepositoryCredentials:  &ecsTypes.RepositoryCredentials{},
		// ResourceRequirements:   []ecsTypes.ResourceRequirement{},
		// Secrets:                []ecsTypes.Secret{},
		// StartTimeout:           new(int32),
		// StopTimeout:            new(int32),
		// SystemControls:         []ecsTypes.SystemControl{},
		// Ulimits:                []ecsTypes.Ulimit{},
		// User:                   new(string),
		// VolumesFrom:            []ecsTypes.VolumeFrom{},
		// WorkingDirectory:       new(string),
	}

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
