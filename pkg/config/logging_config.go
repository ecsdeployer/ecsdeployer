package config

import (
	"errors"
)

// Project level
type LoggingConfig struct {
	Disabled       bool            `yaml:"disabled,omitempty" json:"disabled,omitempty"`
	FirelensConfig *FirelensConfig `yaml:"firelens,omitempty" json:"firelens,omitempty"`
	AwsLogConfig   *AwsLogConfig   `yaml:"awslogs,omitempty" json:"awslogs,omitempty"`
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

	if obj.FirelensConfig.IsDisabled() && obj.AwsLogConfig.IsDisabled() {
		return errors.New("if you want to disable logging, set the 'disabled:true' flag on the 'logging' section")
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

	// disable awslogs if firelens is used
	if !obj.FirelensConfig.IsDisabled() {
		obj.AwsLogConfig.Disabled = true
	}
}

func (obj *LoggingConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type t LoggingConfig // prevent recursive overflow
	var defo = t{}
	if err := unmarshal(&defo); err != nil {
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

/*

func (LoggingConfig) JSONSchemaExtend(base *jsonschema.Schema) {
	base.Description = "Configure logging options"

	// disabledSchema, _ := base.Properties.Get("disabled")
	awsLogsSchema, _ := base.Properties.Get("awslogs")
	firelensSchema, _ := base.Properties.Get("firelens")

	base.Properties = nil

	// falseyDisabledSchema := &jsonschema.Schema{
	// 	Type: "boolean",
	// }

	trueDisabledSchema := &jsonschema.Schema{
		Type: "object",
		Properties: configschema.NewPropertyChain().Set("disabled", &jsonschema.Schema{
			Type:  "boolean",
			Const: true,
		}).End(),
	}

	base.OneOf = []*jsonschema.Schema{
		{
			Type:     "object",
			Required: []string{"disabled"},
			Properties: configschema.NewPropertyChain().Set("disabled", &jsonschema.Schema{
				Type:  "boolean",
				Const: true,
			}).End(),
			AdditionalProperties: jsonschema.FalseSchema,
			Description:          "Entirely disable logging",
		},
		{
			Type:                 "object",
			Required:             []string{"awslogs"},
			Properties:           configschema.NewPropertyChain().Set("awslogs", awsLogsSchema).End(),
			AdditionalProperties: jsonschema.FalseSchema,
			Not:                  trueDisabledSchema,
		},
		{
			Type:                 "object",
			Required:             []string{"firelens"},
			Properties:           configschema.NewPropertyChain().Set("firelens", firelensSchema).End(),
			AdditionalProperties: jsonschema.FalseSchema,
			Not:                  trueDisabledSchema,
		},
	}

}
*/
