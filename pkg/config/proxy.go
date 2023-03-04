package config

import (
	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/invopop/jsonschema"
)

type ProxyConfig struct {
	Disabled bool `yaml:"disabled,omitempty" json:"disabled,omitempty"`

	Type          *string   `yaml:"type,omitempty" json:"type,omitempty"`
	ContainerName *string   `yaml:"container,omitempty" json:"container,omitempty"`
	Properties    EnvVarMap `yaml:"properties,omitempty" json:"properties,omitempty"`
}

func (nc *ProxyConfig) ApplyDefaults() {
	if nc.Type == nil {
		nc.Type = aws.String(string(ecsTypes.ProxyConfigurationTypeAppmesh))
	}

	if nc.ContainerName == nil {
		nc.ContainerName = aws.String("envoy")
	}

	if nc.Properties == nil {
		nc.Properties = make(EnvVarMap)
	}

}

func (nc *ProxyConfig) Validate() error {

	if nc.Disabled {
		return nil
	}

	if util.IsBlank(nc.Type) {
		return NewValidationError("proxy type is required")
	}

	if util.IsBlank(nc.ContainerName) {
		return NewValidationError("proxy container is required")
	}

	if nc.Properties.HasSSM() {
		return NewValidationError("proxy properties cannot reference SSM values")
	}

	return nil
}

func (obj *ProxyConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {

	var boolVal bool
	if err := unmarshal(&boolVal); err == nil {

		if boolVal {
			return NewValidationError("you cannot set a proxy configuration to true, you must specify the parameters.")
		}

		*obj = ProxyConfig{
			Disabled: true,
		}

		return nil

	}

	type tProxyConfig ProxyConfig
	var defo = tProxyConfig{}
	if err := unmarshal(&defo); err != nil {
		return err
	}

	*obj = ProxyConfig(defo)

	obj.ApplyDefaults()

	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

func (ProxyConfig) JSONSchemaExtend(base *jsonschema.Schema) {
	defo := &ProxyConfig{}
	defo.ApplyDefaults()

	configschema.SchemaPropMerge(base, "type", func(s *jsonschema.Schema) {
		s.Default = defo.Type
		s.Examples = util.StrArrayToInterArray(ecsTypes.ProxyConfigurationTypeAppmesh.Values())
		s.Description = "Proxy type. You should omit this unless you know what you are doing."
	})

	configschema.SchemaPropMerge(base, "container", func(s *jsonschema.Schema) {
		s.Default = defo.ContainerName
		s.Description = "Name of the sidecar that provides the proxy"
	})

	orig := *base
	newBase := &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "boolean",
				Description: "Disable proxy configuration for a specific task",
				Const:       false,
			},
			&orig,
		},
	}
	*base = *newBase

}

// public.ecr.aws/appmesh/aws-appmesh-envoy:v1.22.2.1-prod
