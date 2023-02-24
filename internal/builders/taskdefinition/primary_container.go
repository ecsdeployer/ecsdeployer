package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func (b *Builder) applyPrimaryContainer() error {

	b.primaryContainer = &ecsTypes.ContainerDefinition{
		// Command:                []string{},
		// Cpu:                    0,
		// DependsOn:              []ecsTypes.ContainerDependency{},
		// DisableNetworking:      new(bool),
		// DnsSearchDomains:       []string{},
		// DnsServers:             []string{},
		// DockerLabels:           map[string]string{},
		// DockerSecurityOptions:  []string{},
		// EntryPoint:             []string{},
		// Environment:            []ecsTypes.KeyValuePair{},
		// EnvironmentFiles:       []ecsTypes.EnvironmentFile{},
		// Essential:              new(bool),
		// ExtraHosts:             []ecsTypes.HostEntry{},
		// FirelensConfiguration:  &ecsTypes.FirelensConfiguration{},
		// HealthCheck:            &ecsTypes.HealthCheck{},
		// Hostname:               new(string),
		// Image:                  new(string),
		// Interactive:            new(bool),
		// Links:                  []string{},
		// LinuxParameters:        &ecsTypes.LinuxParameters{},
		// LogConfiguration: &ecsTypes.LogConfiguration{
		// 	LogDriver:     "",
		// 	Options:       map[string]string{},
		// 	SecretOptions: []ecsTypes.Secret{},
		// },
		// Memory:                 new(int32),
		// MemoryReservation:      new(int32),
		// MountPoints:            []ecsTypes.MountPoint{},
		// Name:                   new(string),
		// PortMappings:           []ecsTypes.PortMapping{},
		// Privileged:             new(bool),
		// PseudoTerminal:         new(bool),
		// ReadonlyRootFilesystem: new(bool),
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

	if err := b.applyContainerDefaults(b.primaryContainer, b.entity); err != nil {
		return err
	}
	if err := b.addEnvVarsToContainer(b.primaryContainer, b.baseEnvVars); err != nil {
		return err
	}

	if err := b.applyContainerLogging(b.primaryContainer, b.entity); err != nil {
		return err
	}

	if err := b.applyContainerHealthCheck(b.primaryContainer, util.Coalesce(b.commonTask.HealthCheck, b.taskDefaults.HealthCheck)); err != nil {
		return err
	}

	return nil
}
