package config

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/invopop/jsonschema"
)

type FirelensConfig struct {
	Disabled    bool        `yaml:"disabled,omitempty" json:"disabled,omitempty"`
	Type        *string     `yaml:"type,omitempty" json:"type,omitempty"`
	Name        *string     `yaml:"container_name,omitempty" json:"container_name,omitempty"`
	Options     EnvVarMap   `yaml:"options,omitempty" json:"options,omitempty"`
	EnvVars     EnvVarMap   `yaml:"environment,omitempty" json:"environment,omitempty"`
	Credentials *string     `yaml:"credentials,omitempty" json:"credentials,omitempty"`
	InheritEnv  *bool       `yaml:"inherit_env,omitempty" json:"inherit_env,omitempty"`
	Image       *ImageUri   `yaml:"image,omitempty" json:"image,omitempty"`
	Memory      *MemorySpec `yaml:"memory,omitempty" json:"memory,omitempty"`
	// Logging     *TaskLoggingConfig `yaml:"logging,omitempty" json:"logging,omitempty"`

	// should we log the firelens container to AWSLogs (not the app logs, but firelens itself)
	LogToAwsLogs *FirelensAwsLogGroup `yaml:"log_to_awslogs,omitempty" json:"log_to_awslogs,omitempty"`
}

func (obj *FirelensConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var val bool
	if err := unmarshal(&val); err != nil {

		if errors.Is(err, ErrValidation) {
			return err
		}

		type tFirelensConfig FirelensConfig // prevent recursive overflow
		var defo = tFirelensConfig{}
		if err := unmarshal(&defo); err != nil {
			return err
		} else {
			*obj = FirelensConfig(defo)
		}

	} else {
		*obj = FirelensConfig{
			Disabled: !val,
		}
	}

	obj.ApplyDefaults()

	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

func (obj *FirelensConfig) IsDisabled() bool {
	if obj.Disabled {
		return true
	}
	return *obj.Type == LoggingDisableFlag
}

func (obj *FirelensConfig) Validate() error {
	if obj.IsDisabled() {
		return nil
	}

	if obj.Image == nil {
		return NewValidationError("you must provide an image URI for the firelens configuration")
	}

	for _, v := range obj.Options {
		if v.IsSSM() {
			return NewValidationError("you cannot have SSM options in Firelens configuration")
		}
	}

	return nil
}

func (obj *FirelensConfig) ApplyDefaults() {
	if obj.Type == nil {
		obj.Type = aws.String(string(ecsTypes.FirelensConfigurationTypeFluentbit))
	}

	if obj.InheritEnv == nil {
		obj.InheritEnv = aws.Bool(false)
	}

	if obj.LogToAwsLogs == nil {
		obj.LogToAwsLogs = &FirelensAwsLogGroup{}
	}

	if obj.Memory == nil {
		obj.Memory = &MemorySpec{value: 50}
	}

	if obj.Name == nil {
		obj.Name = aws.String("log_router")
	}

	if obj.EnvVars == nil {
		obj.EnvVars = make(EnvVarMap)
	}

	if obj.Options == nil {
		obj.Options = make(EnvVarMap)
	}

	if obj.Image == nil && !obj.IsDisabled() {
		if *obj.Type == string(ecsTypes.FirelensConfigurationTypeFluentbit) {
			obj.Image = &ImageUri{
				uri: aws.String("public.ecr.aws/aws-observability/aws-for-fluent-bit:latest"),
			}
		}
	}

	// if obj.Logging == nil {
	// 	obj.Logging = &TaskLoggingConfig{
	// 		Driver: aws.String(LoggingDisableFlag),
	// 	}
	// }
}

func (FirelensConfig) JSONSchemaExtend(base *jsonschema.Schema) {

	def := &FirelensConfig{}
	def.ApplyDefaults()

	name := configschema.GetSchemaProp(base, "container_name")
	if def.Name != nil {
		name.Default = def.Name
	}

	configschema.SchemaPropMerge(base, "type", func(prop *jsonschema.Schema) {
		if def.Type != nil {
			prop.Default = def.Type
		}
		prop.Enum = []interface{}{
			ecsTypes.FirelensConfigurationTypeFluentbit,
			ecsTypes.FirelensConfigurationTypeFluentd,
		}
	})

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
