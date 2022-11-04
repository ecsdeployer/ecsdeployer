package taskdef

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// TODO: Finish this
func ContainerDefinitionBuilder(ctx *config.Context, common config.CommonContainerAttrs) (*ecsTypes.ContainerDefinition, error) {
	cont := &ecsTypes.ContainerDefinition{
		Command:               []string{},
		Cpu:                   0,
		DependsOn:             []ecsTypes.ContainerDependency{},
		DisableNetworking:     new(bool),
		DockerLabels:          map[string]string{},
		DockerSecurityOptions: []string{},
		EntryPoint:            []string{},
		Environment:           []ecsTypes.KeyValuePair{},
		Image:                 new(string),
		LogConfiguration:      &ecsTypes.LogConfiguration{},
		Memory:                new(int32),
		MemoryReservation:     new(int32),
		Name:                  aws.String(common.Name),
		PortMappings:          []ecsTypes.PortMapping{},
		RepositoryCredentials: &ecsTypes.RepositoryCredentials{},
		Secrets:               []ecsTypes.Secret{},
		StartTimeout:          new(int32),
		StopTimeout:           new(int32),
		User:                  new(string),
		WorkingDirectory:      new(string),
		// Essential:             sidecar.Essential,
		// VolumesFrom:            []ecstypes.VolumeFrom{},
		// Interactive:           new(bool),
		// Links:                 []string{},
		// LinuxParameters:       &ecstypes.LinuxParameters{},
		// DnsSearchDomains:      []string{},
		// DnsServers:            []string{},
		// SystemControls:         []ecstypes.SystemControl{},
		// Ulimits:                []ecstypes.Ulimit{},
		// ResourceRequirements:   []ecstypes.ResourceRequirement{},
		// Privileged:             new(bool),
		// PseudoTerminal:         new(bool),
		// ReadonlyRootFilesystem: new(bool),
		// MountPoints:           []ecstypes.MountPoint{},
		// ExtraHosts:            []ecstypes.HostEntry{},
		// FirelensConfiguration: &ecstypes.FirelensConfiguration{},
		// HealthCheck:           &ecstypes.HealthCheck{},
		// Hostname:              new(string),
		// EnvironmentFiles:      []ecstypes.EnvironmentFile{},
	}

	return cont, nil
}
