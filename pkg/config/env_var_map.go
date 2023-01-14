package config

import (
	"github.com/invopop/jsonschema"
	"golang.org/x/exp/maps"
)

type EnvVarMap map[string]EnvVar

func (EnvVarMap) JSONSchemaExtend(base *jsonschema.Schema) {

	patt := base.PatternProperties
	patt["^[a-zA-Z_][^=]*$"] = patt[".*"]
	delete(patt, ".*")
	base.AdditionalProperties = jsonschema.FalseSchema
}

// Filters a map of env vars and removes any Unset values
func (obj EnvVarMap) Filter() EnvVarMap {
	newMap := make(EnvVarMap, len(obj))
	for k, v := range obj {
		if v.IsUnset() {
			continue
		}
		v := v
		newMap[k] = v
	}
	return newMap
}

// are any of the values inside an SSM reference
func (obj EnvVarMap) HasSSM() bool {
	for _, v := range obj {
		if v.IsSSM() {
			return true
		}
	}
	return false
}

func MergeEnvVarMaps(values ...EnvVarMap) EnvVarMap {
	newMap := make(EnvVarMap)
	for _, value := range values {
		value := value
		maps.Copy(newMap, value)
	}

	return newMap
}
