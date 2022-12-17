package testutil

import (
	"fmt"
	"strings"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"github.com/xeipuuv/gojsonschema"
)

type SchemaChecker struct {
	schema *gojsonschema.Schema
}

func NewSchemaChecker(v any) *SchemaChecker {
	schema := configschema.GenerateSchema(v)

	schemaJson, err := util.Jsonify(schema)
	if err != nil {
		panic(err)
	}

	jsonSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(schemaJson))
	if err != nil {
		panic(err)
	}

	return &SchemaChecker{
		schema: jsonSchema,
	}
}

// Validates that the YAML provided adheres to the Schema (when converted to json)
func (obj *SchemaChecker) CheckYAML(t *testing.T, val string) error {
	t.Helper()

	var tmp interface{}

	if err := yaml.Unmarshal([]byte(val), &tmp); err != nil {
		panic(err)
	}

	jsonData, err := util.Jsonify(tmp)
	if err != nil {
		panic(err)
	}

	return obj.CheckJSON(t, string(jsonData))

}

// Validates that the JSON provided adheres to the Schema
func (obj *SchemaChecker) CheckJSON(t *testing.T, val string) error {
	t.Helper()

	valLoader := gojsonschema.NewStringLoader(val)
	result, err := obj.schema.Validate(valLoader)
	if err != nil {
		panic(err)
	}

	if result.Valid() {
		return nil
	}

	schemaErrs := result.Errors()

	errList := make([]string, 0, len(schemaErrs)+1)

	errList = append(errList, "SchemaFailure")

	for _, err := range schemaErrs {
		errList = append(errList, err.String())
	}
	return fmt.Errorf(strings.Join(errList, "\n"))
}
