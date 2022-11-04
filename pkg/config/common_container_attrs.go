package config

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
}

func (cta *CommonContainerAttrs) Validate() error {

	return nil
}

func (cta *CommonContainerAttrs) TemplateFields() map[string]interface{} {
	return map[string]interface{}{
		"Name": cta.Name,
	}
}
