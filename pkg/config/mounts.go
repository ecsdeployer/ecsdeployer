package config

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type Mount struct {
	ContainerPath string `yaml:"path" json:"path"`
	SourceVolume  string `yaml:"source" json:"source"`
	ReadOnly      bool   `yaml:"readonly,omitempty" json:"readonly,omitempty"`
}

func (obj *Mount) Validate() error {

	if util.IsBlank(&obj.ContainerPath) {
		return NewValidationError("mount path cannot be empty")
	}

	if util.IsBlank(&obj.SourceVolume) {
		return NewValidationError("mount source cannot be empty")
	}

	return nil
}

func (obj *Mount) ApplyDefaults() {
}

func (obj *Mount) ToAws() ecsTypes.MountPoint {
	return ecsTypes.MountPoint{
		ContainerPath: aws.String(obj.ContainerPath),
		SourceVolume:  aws.String(obj.SourceVolume),
		ReadOnly:      aws.Bool(obj.ReadOnly),
	}
}

func (obj *Mount) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tMount Mount
	var defo = tMount{}
	if err := unmarshal(&defo); err != nil {
		return err
	}

	*obj = Mount(defo)

	obj.ApplyDefaults()

	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}
