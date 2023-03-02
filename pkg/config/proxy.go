package config

import (
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type AppMesh struct {
	Type          *string           `yaml:"type,omitempty" json:"type,omitempty"`
	ContainerName *string           `yaml:"container_name,omitempty" json:"container_name,omitempty"`
	Properties    map[string]string `yaml:"properties,omitempty" json:"properties,omitempty"`
}

func (nc *AppMesh) ApplyDefaults() {
}

func (nc *AppMesh) Validate() error {

	_ = ecsTypes.ProxyConfigurationTypeAppmesh

	return nil
}

// public.ecr.aws/appmesh/aws-appmesh-envoy:v1.22.2.1-prod
