package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type Container struct {
	ParentTask *Task

	IsPrimary bool

	Name       string
	Command    *config.ShellCommand
	EntryPoint *config.ShellCommand
	Image      config.ImageUri
	Labels     map[string]string
	DependsOn  map[string]ecsTypes.ContainerCondition

	Environment config.EnvVarMap

	Cpu    *config.CpuSpec
	Memory *config.MemorySpec

	Credentials *string

	StartTimeout *config.Duration
	StopTimeout  *config.Duration

	HealthCheck *config.HealthCheck

	Firelens  *FirelensConfig
	LogConfig *LogConfig
}

func (cd *Container) ctx() *config.Context {
	return cd.ParentTask.Context
}

func (cd *Container) EvalTpl(tplStr string) (string, error) {
	return cd.ParentTask.EvalTpl(tplStr)
}

func (cd *Container) Export() (*ecsTypes.ContainerDefinition, error) {
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
