package config

import (
	"maps"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/invopop/jsonschema"
)

type EnvVarMap map[string]EnvVar

func (EnvVarMap) JSONSchemaExtend(base *jsonschema.Schema) {

	// patt := base.PatternProperties
	// if len(patt) > 0 {
	// 	patt["^[a-zA-Z_][^=]*$"] = patt[".*"]
	// 	delete(patt, ".*")
	// } else {
	// 	base.PatternProperties = map[string]*jsonschema.Schema{
	// 		"^[a-zA-Z_][^=]*$": {Ref: "#/$defs/EnvVar"},
	// 	}
	// }
	base.PatternProperties = map[string]*jsonschema.Schema{
		"^[a-zA-Z_][^=]*$": {Ref: "#/$defs/EnvVar"},
	}
	base.AdditionalProperties = jsonschema.FalseSchema
}

// func (EnvVarMap) JSONSchema() *jsonschema.Schema {
// 	return &jsonschema.Schema{
// 		Type:                 "object",
// 		AdditionalProperties: jsonschema.FalseSchema,
// 		PatternProperties: map[string]*jsonschema.Schema{
// 			"^[a-zA-Z_][^=]*$": {Ref: "#/$defs/EnvVar"},
// 		},
// 	}
// }

// Filters a map of env vars and removes any Unset values
// You should only use this at the very very end of an evaluation tree.
// (i.e. after merging parent maps)
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
		maps.Copy(newMap, value)
	}

	return newMap
}

type EnvVarExportFunc[T any] func(string, string) T

func ExportEnvVarMap[Te any, Ts any](varMap EnvVarMap, tpl templater, envVarFunc EnvVarExportFunc[Te], secretFunc EnvVarExportFunc[Ts]) ([]Te, []Ts, error) {
	var envvars = []Te{}
	var secrets = []Ts{}

	for key, val := range varMap.Filter() {
		if val.IsSSM() {

			secrets = append(secrets, secretFunc(key, util.Must(val.GetValue(nil))))
			continue
		}

		value, err := val.GetValue(tpl)
		if err != nil {
			return nil, nil, err
		}

		if util.IsBlank(&value) {
			continue
		}

		envvars = append(envvars, envVarFunc(key, value))
	}

	return envvars, secrets, nil
}
