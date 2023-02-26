package config

import (
	"errors"

	"github.com/invopop/jsonschema"
)

// Project level
type LoggingConfig struct {
	Disabled       bool             `yaml:"disabled,omitempty" json:"disabled,omitempty"`
	FirelensConfig *FirelensConfig  `yaml:"firelens,omitempty" json:"firelens,omitempty"`
	AwsLogConfig   *AwsLogConfig    `yaml:"awslogs,omitempty" json:"awslogs,omitempty"`
	Custom         *CustomLogConfig `yaml:"custom,omitempty" json:"custom,omitempty"`
}

func (obj *LoggingConfig) IsDisabled() bool {
	return obj.Disabled
}

func (obj *LoggingConfig) Validate() error {
	if obj.IsDisabled() {
		return nil
	}

	if err := obj.AwsLogConfig.Validate(); err != nil {
		return err
	}

	if err := obj.FirelensConfig.Validate(); err != nil {
		return err
	}

	if err := obj.Custom.Validate(); err != nil {
		return err
	}

	if obj.FirelensConfig.IsDisabled() && obj.AwsLogConfig.IsDisabled() && obj.Custom.IsDisabled() {
		return NewValidationError("if you want to disable logging, set the 'disabled:true' flag on the 'logging' section")
	}

	return nil
}

func (obj *LoggingConfig) ApplyDefaults() {
	if obj.FirelensConfig == nil {
		obj.FirelensConfig = &FirelensConfig{
			Disabled: true,
		}
	}
	obj.FirelensConfig.ApplyDefaults()

	if obj.AwsLogConfig == nil {
		obj.AwsLogConfig = &AwsLogConfig{}
	}
	obj.AwsLogConfig.ApplyDefaults()

	if obj.Custom == nil {
		obj.Custom = &CustomLogConfig{}
	}
	obj.Custom.ApplyDefaults()

	// disable awslogs if firelens is used
	if !obj.FirelensConfig.IsDisabled() {
		obj.AwsLogConfig.Disabled = true
	}

	if !obj.Custom.IsDisabled() {
		obj.FirelensConfig.Disabled = true
		obj.AwsLogConfig.Disabled = true
	}
}

func (obj *LoggingConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tLoggingConfig LoggingConfig // prevent recursive overflow
	var defo = tLoggingConfig{}
	if err := unmarshal(&defo); err != nil {

		if errors.Is(err, ErrValidation) {
			return err
		}

		var val bool
		if err := unmarshal(&val); err != nil {
			return err
		}
		*obj = LoggingConfig{
			Disabled: !val,
		}
	} else {
		*obj = LoggingConfig(defo)
	}

	obj.ApplyDefaults()

	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

func (LoggingConfig) JSONSchemaExtend(base *jsonschema.Schema) {
	orig := *base
	newBase := &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type: "boolean",
			},
			&orig,
		},
	}
	*base = *newBase
}
