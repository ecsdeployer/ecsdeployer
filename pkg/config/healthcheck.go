package config

import (
	"github.com/invopop/jsonschema"
)

type HealthCheck struct {
	Disabled bool `yaml:"disabled,omitempty" json:"disabled,omitempty"`

	Command     ShellCommand `yaml:"command,omitempty" json:"command,omitempty"`
	Interval    *Duration    `yaml:"interval,omitempty" json:"interval,omitempty"`
	Retries     *int32       `yaml:"retries,omitempty" json:"retries,omitempty" jsonschema_extras:"minimum=1"`
	StartPeriod *Duration    `yaml:"start_period,omitempty" json:"start_period,omitempty"`
	Timeout     *Duration    `yaml:"timeout,omitempty" json:"timeout,omitempty"`
}

func (obj *HealthCheck) Validate() error {

	if obj.Disabled {
		return nil
	}

	if obj.Command == nil || len(obj.Command) == 0 {
		return NewValidationError("Healthcheck command cannot be empty")
	}

	if obj.Command[0] != "CMD" && obj.Command[0] != "CMD-SHELL" {
		return NewValidationError("Healthcheck command MUST start with 'CMD' or 'CMD-SHELL'")
	}

	if obj.Retries != nil && *obj.Retries < 0 {
		return NewValidationError("Retries cannot be negative")
	}

	return nil
}

func (obj *HealthCheck) ApplyDefaults() {
}

func (obj *HealthCheck) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var boolVal bool
	if err := unmarshal(&boolVal); err == nil {

		if boolVal {
			return NewValidationError("you cannot set a health check to true, you must specify the parameters.")
		}

		*obj = HealthCheck{
			Disabled: true,
		}

		return nil

	}

	type tHealthCheck HealthCheck
	var defo = tHealthCheck{}
	if err := unmarshal(&defo); err != nil {
		return err
	}

	*obj = HealthCheck(defo)

	obj.ApplyDefaults()

	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

func (HealthCheck) JSONSchemaExtend(base *jsonschema.Schema) {

	defo := &HealthCheck{}
	defo.ApplyDefaults()

	orig := *base
	newBase := &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "boolean",
				Description: "Disable a healthcheck for a specific task",
				Const:       false,
			},
			&orig,
		},
	}
	*base = *newBase
}
