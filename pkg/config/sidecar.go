package config

import (
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
)

type Sidecar struct {
	CommonContainerAttrs `yaml:",inline" json:",inline"`

	InheritEnv        bool          `yaml:"inherit_env,omitempty" json:"inherit_env,omitempty"`
	PortMappings      []PortMapping `yaml:"port_mappings,omitempty" json:"port_mappings,omitempty"`
	MemoryReservation *MemorySpec   `yaml:"memory_reservation,omitempty" json:"memory_reservation,omitempty"`
	Essential         *bool         `yaml:"essential,omitempty" json:"essential,omitempty"`
}

func (obj *Sidecar) Validate() error {

	if obj.Name == "" {
		return errors.New("you must set a name")
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
