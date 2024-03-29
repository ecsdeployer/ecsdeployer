package config

import (
	"encoding/json"
	"errors"

	"github.com/invopop/jsonschema"
)

type LoggingType uint8

const (
	LoggingTypeDisabled LoggingType = iota
	LoggingTypeAwslogs
	LoggingTypeFirelens
	LoggingTypeCustom
)

// Project level
type LoggingConfig struct {
	Disabled       bool             `yaml:"disabled,omitempty" json:"disabled,omitempty"`
	FirelensConfig *FirelensConfig  `yaml:"firelens,omitempty" json:"firelens,omitempty"`
	AwsLogConfig   *AwsLogConfig    `yaml:"awslogs,omitempty" json:"awslogs,omitempty"`
	Custom         *CustomLogConfig `yaml:"custom,omitempty" json:"custom,omitempty"`
}

func (obj *LoggingConfig) Type() LoggingType {

	if obj.IsDisabled() {
		return LoggingTypeDisabled
	}

	if !obj.Custom.IsDisabled() {
		return LoggingTypeCustom
	}

	if !obj.FirelensConfig.IsDisabled() {
		return LoggingTypeFirelens
	}

	if !obj.AwsLogConfig.IsDisabled() {
		return LoggingTypeAwslogs
	}

	return LoggingTypeDisabled
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

	// if obj.FirelensConfig.IsDisabled() && obj.AwsLogConfig.IsDisabled() && obj.Custom.IsDisabled() {
	// 	return NewValidationError("if you want to disable logging, set the 'disabled:true' flag on the 'logging' section")
	// }

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
	// if !obj.FirelensConfig.IsDisabled() {
	// 	obj.AwsLogConfig.Disabled = true
	// }

	// if !obj.Custom.IsDisabled() {
	// 	obj.FirelensConfig.Disabled = true
	// 	obj.AwsLogConfig.Disabled = true
	// }
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

func (obj *LoggingConfig) MarshalYAML() (interface{}, error) {
	if obj.IsDisabled() {
		return false, nil
	}
	type tLoggingConfig LoggingConfig

	switch obj.Type() {
	case LoggingTypeAwslogs:
		return &tLoggingConfig{AwsLogConfig: obj.AwsLogConfig}, nil
	case LoggingTypeFirelens:
		return &tLoggingConfig{FirelensConfig: obj.FirelensConfig}, nil
	case LoggingTypeCustom:
		return &tLoggingConfig{Custom: obj.Custom}, nil
	default:
		return "!!!!!!!!!UNKNOWN", nil
	}
}

func (obj *LoggingConfig) MarshalJSON() ([]byte, error) {
	data, err := obj.MarshalYAML()
	if err != nil {
		return nil, err
	}
	return json.Marshal(data)
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
