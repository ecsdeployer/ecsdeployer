package config

import (
	"strings"

	"github.com/invopop/jsonschema"
)

type KeepInSyncTaskDefinitions uint8

const (
	_ KeepInSyncTaskDefinitions = iota

	// dont manage task definitons at all
	KeepInSyncTaskDefinitionsDisabled

	// deregister all old task definitions (even if we no longer manage it)
	KeepInSyncTaskDefinitionsEnabled

	// only deregister the previous task definition for things we are actively managing
	KeepInSyncTaskDefinitionsOnlyManaged
)

const (
	kisTaskDefOnlyManagedStr = "ONLY_MANAGED"
)

func (obj *KeepInSyncTaskDefinitions) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var boolVal bool
	if err := unmarshal(&boolVal); err == nil {
		if boolVal {
			*obj = KeepInSyncTaskDefinitionsEnabled
		} else {
			*obj = KeepInSyncTaskDefinitionsDisabled
		}
		return nil
	}

	var strVal string
	if err := unmarshal(&strVal); err == nil {
		switch strings.ToUpper(strVal) {
		case kisTaskDefOnlyManagedStr:
			*obj = KeepInSyncTaskDefinitionsOnlyManaged
			return nil
		case "TRUE":
			*obj = KeepInSyncTaskDefinitionsEnabled
			return nil
		case "FALSE":
			*obj = KeepInSyncTaskDefinitionsDisabled
			return nil
		}

	}

	return NewValidationError("Invalid value for keep_in_sync.task_definitions. Must be one of true, false, only_managed")

}

func (KeepInSyncTaskDefinitions) JSONSchema() *jsonschema.Schema {
	return &jsonschema.Schema{
		Description: "How to keep old task definitions in sync",
		OneOf: []*jsonschema.Schema{
			{
				Type:        "string",
				Const:       kisTaskDefOnlyManagedStr,
				Description: "Will only deregister previous definitions for tasks defined in project. Will not manage definitions for removed tasks.",
			},
			{
				Type:        "boolean",
				Const:       true,
				Description: "Enables keeping task definitions in sync. Previous definitions as well as removed tasks will be synced.",
			},
			{
				Type:        "boolean",
				Const:       false,
				Description: "Disables keeping all task definitions in sync. Previous definitions will not be deregistered.",
			},
		},
		Default: true,
	}
}
