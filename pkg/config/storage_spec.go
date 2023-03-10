package config

import (
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/invopop/jsonschema"
)

type StorageSpec int32

func (ss StorageSpec) Gigabytes() int32 {
	return int32(ss)
}

func (ss StorageSpec) ToAws() *ecsTypes.EphemeralStorage {
	return &ecsTypes.EphemeralStorage{
		SizeInGiB: ss.Gigabytes(),
	}
}

func NewStorageSpec(gb int32) (*StorageSpec, error) {
	spec := StorageSpec(gb)
	return &spec, nil
}

func (StorageSpec) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "integer",
		Description: "Size in GB of storage to add",
		Minimum:     20,
		Default:     20,
	}
}
