package config

import (
	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/invopop/jsonschema"
)

type ProxyConfig struct {
	Type          *string   `yaml:"type,omitempty" json:"type,omitempty"`
	ContainerName *string   `yaml:"container_name,omitempty" json:"container_name,omitempty"`
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

	if util.IsBlank(nc.Type) {
		return NewValidationError("proxy type is required")
	}

	if util.IsBlank(nc.ContainerName) {
		return NewValidationError("proxy container_name is required")
	}

	if nc.Properties.HasSSM() {
		return NewValidationError("proxy properties cannot reference SSM values")
	}

	return nil
}

func (obj *ProxyConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
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

	configschema.SchemaPropMerge(base, "container_name", func(s *jsonschema.Schema) {
		s.Default = defo.ContainerName
		s.Description = "Name of the sidecar that provides the proxy"
	})

}

// public.ecr.aws/appmesh/aws-appmesh-envoy:v1.22.2.1-prod
