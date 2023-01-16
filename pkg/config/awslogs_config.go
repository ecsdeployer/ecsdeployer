package config

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/invopop/jsonschema"
)

// Project level
type AwsLogConfig struct {
	Disabled  bool          `yaml:"disabled,omitempty" json:"disabled,omitempty"`
	Retention *LogRetention `yaml:"retention,omitempty" json:"retention,omitempty"`
	Options   EnvVarMap     `yaml:"options,omitempty" json:"options,omitempty"`
}

func (obj *AwsLogConfig) IsDisabled() bool {
	return obj.Disabled
}

func (obj *AwsLogConfig) Validate() error {
	if obj.IsDisabled() {
		return nil
	}

	return nil
}

func (obj *AwsLogConfig) ApplyDefaults() {
	if obj.Retention == nil {
		// obj.Retention = aws.Int32(365)
		obj.Retention = util.Ptr(util.Must(ParseLogRetention(defaultLogRetention)))
	}

	if obj.Options == nil {
		obj.Options = make(EnvVarMap)
	}
}

func (obj *AwsLogConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tAwsLogConfig AwsLogConfig
	var defo = tAwsLogConfig{}
	if err := unmarshal(&defo); err != nil {

		if errors.Is(err, ErrValidation) {
			return err
		}

		var val bool
		if err := unmarshal(&val); err != nil {
			return err
		}
		*obj = AwsLogConfig{
			Disabled: !val,
		}
	} else {
		*obj = AwsLogConfig(defo)
	}

	obj.ApplyDefaults()

	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

func (AwsLogConfig) JSONSchemaExtend(base *jsonschema.Schema) {
	orig := *base
	newBase := &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "boolean",
				Description: "Enable or disable logging to AWS using default values",
			},
			&orig,
		},
	}
	*base = *newBase
}
