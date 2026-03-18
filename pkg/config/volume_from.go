package config

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/internal/util"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/invopop/jsonschema"
)

type VolumeFrom struct {
	SourceContainer string `yaml:"source" json:"source"`
	ReadOnly        bool   `yaml:"readonly,omitempty" json:"readonly,omitempty"`
}

func (obj *VolumeFrom) Validate() error {

	if util.IsBlank(&obj.SourceContainer) {
		return NewValidationError("source container cannot be empty")
	}

	return nil
}

func (obj *VolumeFrom) ApplyDefaults() {
}

func (obj *VolumeFrom) ToAws() ecsTypes.VolumeFrom {
	return ecsTypes.VolumeFrom{
		SourceContainer: new(obj.SourceContainer),
		ReadOnly:        new(obj.ReadOnly),
	}
}

func (obj *VolumeFrom) UnmarshalYAML(unmarshal func(any) error) error {
	type tVolumeFrom VolumeFrom
	var defo = tVolumeFrom{}
	if err := unmarshal(&defo); err != nil {

		if errors.Is(err, ErrValidation) {
			return err
		}

		var str string
		if err := unmarshal(&str); err != nil {
			return err
		}

		*obj = VolumeFrom{SourceContainer: str}
	} else {
		*obj = VolumeFrom(defo)
	}

	obj.ApplyDefaults()

	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

func (VolumeFrom) JSONSchemaExtend(base *jsonschema.Schema) {
	orig := *base
	newBase := &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{
				Type:        "string",
				Description: "Shorthand to specify a container name",
			},
			&orig,
		},
	}
	*base = *newBase
}
