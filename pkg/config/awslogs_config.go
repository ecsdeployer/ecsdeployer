package config

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
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
	type t AwsLogConfig
	var defo = t{}
	if err := unmarshal(&defo); err != nil {
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
