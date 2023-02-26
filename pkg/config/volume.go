package config

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/invopop/jsonschema"
)

type Volume struct {
	Name      string           `yaml:"name" json:"name"`
	EFSConfig *VolumeEFSConfig `yaml:"efs,omitempty" json:"efs,omitempty"`
}

func (obj *Volume) Validate() error {

	if util.IsBlank(&obj.Name) {
		return NewValidationError("volume name cannot be empty")
	}

	if obj.EFSConfig != nil {
		if err := obj.EFSConfig.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (obj *Volume) ApplyDefaults() {
}

func (obj *Volume) ToAws() ecsTypes.Volume {
	vol := ecsTypes.Volume{
		Name: aws.String(obj.Name),
	}

	if obj.EFSConfig != nil {
		vol.EfsVolumeConfiguration = obj.EFSConfig.ToAws()
	}

	return vol
}

func (obj *Volume) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tVolume Volume
	var defo = tVolume{}
	if err := unmarshal(&defo); err != nil {

		if errors.Is(err, ErrValidation) {
			return err
		}

		var str string
		if err := unmarshal(&str); err != nil {
			return err
		}

		*obj = Volume{Name: str}
	} else {
		*obj = Volume(defo)
	}

	obj.ApplyDefaults()

	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

func (Volume) JSONSchemaExtend(base *jsonschema.Schema) {
	orig := *base
	newBase := &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "string",
				Description: "Shorthand to specify a default bind volume",
			},
			&orig,
		},
	}
	*base = *newBase
}
