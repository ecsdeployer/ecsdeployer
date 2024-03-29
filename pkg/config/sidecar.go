package config

import (
	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/invopop/jsonschema"
)

type Sidecar struct {
	CommonContainerAttrs `yaml:",inline" json:",inline"`

	InheritEnv        bool          `yaml:"inherit_env,omitempty" json:"inherit_env,omitempty"`
	PortMappings      []PortMapping `yaml:"port_mappings,omitempty" json:"port_mappings,omitempty"`
	MemoryReservation *MemorySpec   `yaml:"memory_reservation,omitempty" json:"memory_reservation,omitempty"`
	Essential         *bool         `yaml:"essential,omitempty" json:"essential,omitempty"`
}

func (obj *Sidecar) GetCommonContainerAttrs() CommonContainerAttrs {
	return obj.CommonContainerAttrs
}

func (obj *Sidecar) Validate() error {

	if obj.Name == "" {
		return NewValidationError("you must set a name")
	}

	if err := obj.CommonContainerAttrs.Validate(); err != nil {
		return err
	}

	if len(obj.PortMappings) > 0 {
		for _, portmap := range obj.PortMappings {
			if err := portmap.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (obj *Sidecar) ApplyDefaults() {
	if obj.Essential == nil {
		obj.Essential = aws.Bool(true)
	}

}

func (obj *Sidecar) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type t Sidecar
	var defo = t{}
	if err := unmarshal(&defo); err != nil {
		return err
	}

	*obj = Sidecar(defo)

	obj.ApplyDefaults()

	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

func (Sidecar) JSONSchemaExtend(base *jsonschema.Schema) {

	def := &Sidecar{}
	def.ApplyDefaults()

	configschema.SchemaPropMerge(base, "essential", func(s *jsonschema.Schema) {
		s.Default = def.Essential
	})

	configschema.SchemaPropMerge(base, "inherit_env", func(s *jsonschema.Schema) {
		s.Default = def.InheritEnv
	})

	if base.Required == nil {
		base.Required = make([]string, 0)
	}
	base.Required = append(base.Required, "name")
}
