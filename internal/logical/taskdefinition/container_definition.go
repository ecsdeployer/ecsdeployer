package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type ContainerDefinition struct {
	ParentTaskDef *TaskDefinition

	Name       string
	Command    *config.ShellCommand
	EntryPoint *config.ShellCommand
	Image      config.ImageUri
	Labels     map[string]string

	Cpu    *config.CpuSpec
	Memory *config.MemorySpec

	Credentials *string

	HealthCheck *config.HealthCheck
}

func (td *ContainerDefinition) Export() (*ecsTypes.ContainerDefinition, error) {
	cdef := &ecsTypes.ContainerDefinition{
		Command:               []string{},
		Cpu:                   0,
		DependsOn:             []ecsTypes.ContainerDependency{},
		DockerLabels:          map[string]string{},
		EntryPoint:            []string{},
		Environment:           []ecsTypes.KeyValuePair{},
		Essential:             new(bool),
		FirelensConfiguration: &ecsTypes.FirelensConfiguration{},
		HealthCheck:           &ecsTypes.HealthCheck{},
		Image:                 new(string),
		LinuxParameters:       &ecsTypes.LinuxParameters{},
		LogConfiguration:      &ecsTypes.LogConfiguration{},
		Memory:                new(int32),
		MemoryReservation:     new(int32),
		Name:                  new(string),
		PortMappings:          []ecsTypes.PortMapping{},
		RepositoryCredentials: &ecsTypes.RepositoryCredentials{},
		Secrets:               []ecsTypes.Secret{},
		StartTimeout:          new(int32),
		StopTimeout:           new(int32),
	}
	return cdef, nil
}
