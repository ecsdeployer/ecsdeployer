package config

import (
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/invopop/jsonschema"
	"golang.org/x/exp/maps"
)

type VolumeList map[string]Volume

func (list VolumeList) Validate() error {

	for _, vol := range list {
		if err := vol.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (vlist VolumeList) ToAws() []ecsTypes.Volume {
	out := make([]ecsTypes.Volume, 0)

	for _, val := range vlist {
		out = append(out, val.ToAws())
	}

	return out
}

func (obj *VolumeList) UnmarshalYAML(unmarshal func(interface{}) error) error {
	type tVolumeList []Volume
	var defo = tVolumeList{}
	if err := unmarshal(&defo); err != nil {
		return err
	}

	newMap := make(VolumeList)

	for _, val := range defo {
		if _, ok := newMap[val.Name]; ok {
			return NewValidationError("Duplicate volume name: %s", val.Name)
		}
		newMap[val.Name] = val
	}
	*obj = newMap

	if err := obj.Validate(); err != nil {
		return err
	}

	return nil
}

func (VolumeList) JSONSchemaExtend(base *jsonschema.Schema) {

	// get the original mapping
	volRef := base.PatternProperties[".*"]

	*base = jsonschema.Schema{
		Type:  "array",
		Items: volRef,
	}

}

func MergeVolumeLists(volumeLists ...VolumeList) VolumeList {
	newMap := make(VolumeList)
	for _, value := range volumeLists {
		if value == nil {
			continue
		}
		value := value
		maps.Copy(newMap, value)
	}

	return newMap
}
