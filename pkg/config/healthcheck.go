package config

import "errors"

type HealthCheck struct {
	Command     []string  `yaml:"command" json:"command"`
	Interval    *Duration `yaml:"interval,omitempty" json:"interval,omitempty"`
	Retries     *int32    `yaml:"retries,omitempty" json:"retries,omitempty"`
	StartPeriod *Duration `yaml:"start_period,omitempty" json:"start_period,omitempty"`
	Timeout     *Duration `yaml:"timeout,omitempty" json:"timeout,omitempty"`
}

func (obj *HealthCheck) Validate() error {

	if obj.Command == nil || len(obj.Command) == 0 {
		return errors.New("Healthcheck command cannot be empty")
	}

	if obj.Command[0] != "CMD" && obj.Command[0] != "CMD-SHELL" {
		return errors.New("Healthcheck command MUST start with 'CMD' or 'CMD-SHELL'")
	}

	if obj.Retries != nil && *obj.Retries < 0 {
		return errors.New("Retries cannot be negative")
	}

	return nil
}

func (obj *HealthCheck) ApplyDefaults() {
}
