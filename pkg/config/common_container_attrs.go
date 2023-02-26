package config

import "ecsdeployer.com/ecsdeployer/internal/util"

type CommonContainerAttrs struct {
	Name          string             `yaml:"name,omitempty" json:"name,omitempty" jsonschema:"pattern=^[a-zA-Z][-_a-zA-Z0-9]*$"`
	Command       *ShellCommand      `yaml:"command,omitempty" json:"command,omitempty"`
	EntryPoint    *ShellCommand      `yaml:"entrypoint,omitempty" json:"entrypoint,omitempty"`
	Image         *ImageUri          `yaml:"image,omitempty" json:"image,omitempty"`
	Credentials   *string            `yaml:"credentials,omitempty" json:"credentials,omitempty"`
	Cpu           *CpuSpec           `yaml:"cpu,omitempty" json:"cpu,omitempty"`
	Memory        *MemorySpec        `yaml:"memory,omitempty" json:"memory,omitempty"`
	EnvVars       EnvVarMap          `yaml:"environment,omitempty" json:"environment,omitempty"`
	StartTimeout  *Duration          `yaml:"start_timeout,omitempty" json:"start_timeout,omitempty"`
	StopTimeout   *Duration          `yaml:"stop_timeout,omitempty" json:"stop_timeout,omitempty"`
	DockerLabels  []NameValuePair    `yaml:"labels,omitempty" json:"labels,omitempty"`
	DependsOn     []DependsOn        `yaml:"depends_on,omitempty" json:"depends_on,omitempty"`
	LoggingConfig *TaskLoggingConfig `yaml:"logging,omitempty" json:"logging,omitempty"`
	HealthCheck   *HealthCheck       `yaml:"healthcheck,omitempty" json:"healthcheck,omitempty"`
	MountPoints   []Mount            `yaml:"mounts,omitempty" json:"mounts,omitempty"`
	Ulimits       []Ulimit           `yaml:"ulimits,omitempty" json:"ulimits,omitempty"`
	User          *string            `yaml:"user,omitempty" json:"user,omitempty"`
	Workdir       *string            `yaml:"workdir,omitempty" json:"workdir,omitempty"`
	VolumesFrom   []VolumeFrom       `yaml:"volumes_from,omitempty" json:"volumes_from,omitempty"`
}

func (obj *CommonContainerAttrs) GetCommonContainerAttrs() CommonContainerAttrs {
	return *obj
}

func (cta *CommonContainerAttrs) Validate() error {

	return nil
}

func (cta *CommonContainerAttrs) TemplateFields() map[string]interface{} {
	return map[string]interface{}{
		"Name": cta.Name,
	}
}

func (cta CommonContainerAttrs) NewDefaultedBy(other CommonContainerAttrs) CommonContainerAttrs {
	newCta := CommonContainerAttrs{
		Name:          cta.Name,
		Command:       util.Coalesce(cta.Command, other.Command),
		EntryPoint:    util.Coalesce(cta.EntryPoint, other.EntryPoint),
		Image:         util.Coalesce(cta.Image, other.Image),
		Credentials:   util.Coalesce(cta.Credentials, other.Credentials),
		Cpu:           util.Coalesce(cta.Cpu, other.Cpu),
		Memory:        util.Coalesce(cta.Memory, other.Memory),
		EnvVars:       MergeEnvVarMaps(other.EnvVars, cta.EnvVars),
		StartTimeout:  util.Coalesce(cta.StartTimeout, other.StartTimeout),
		StopTimeout:   util.Coalesce(cta.StopTimeout, other.StopTimeout),
		DockerLabels:  util.CoalesceArray(cta.DockerLabels, other.DockerLabels),
		DependsOn:     util.CoalesceArray(cta.DependsOn, other.DependsOn),
		LoggingConfig: util.Coalesce(cta.LoggingConfig, other.LoggingConfig),
		HealthCheck:   util.Coalesce(cta.HealthCheck, other.HealthCheck),
		Workdir:       util.Coalesce(cta.Workdir, other.Workdir),
		User:          util.Coalesce(cta.User, other.User),
		Ulimits:       util.CoalesceArray(cta.Ulimits, other.Ulimits),
		MountPoints:   util.CoalesceArray(cta.MountPoints, other.MountPoints),
		VolumesFrom:   util.CoalesceArray(cta.VolumesFrom, other.VolumesFrom),
	}

	return newCta
}
