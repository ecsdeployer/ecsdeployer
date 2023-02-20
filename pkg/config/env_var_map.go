package config

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
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
		value := value
		maps.Copy(newMap, value)
	}

	return newMap
}

type EnvVarExportFunc func(string, string) any

func (obj EnvVarMap) Export(tpl templater, envVarFunc EnvVarExportFunc, secretFunc EnvVarExportFunc) (any, any, error) {
	var envvars = []any{}
	var secrets = []any{}

	for key, val := range obj.Filter() {
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
