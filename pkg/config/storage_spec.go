package config

import "github.com/invopop/jsonschema"

type StorageSpec int32

func (ss StorageSpec) Gigabytes() int32 {
	return int32(ss)
}

func (StorageSpec) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Type:        "integer",
		Description: "Size in GB of storage to add",
		Minimum:     20,
		Default:     20,
	}
}
