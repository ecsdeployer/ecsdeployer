package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"github.com/xeipuuv/gojsonschema"
)

type SchemaTester[T interface{}] struct {
	schema *gojsonschema.Schema
	tst    *testing.T
}

func NewSchemaTester[Ttype interface{}](t *testing.T, v interface{}) *SchemaTester[Ttype] {

	t.Helper()
	schema := configschema.GenerateSchema(v)

	schemaJson, err := util.Jsonify(schema)
	if err != nil {
		panic(err)
	}

	jsonSchema, err := gojsonschema.NewSchema(gojsonschema.NewStringLoader(schemaJson))
	if err != nil {
		panic(err)
	}

	return &SchemaTester[Ttype]{
		tst:    t,
		schema: jsonSchema,
	}

}

func (st *SchemaTester[T]) AssertValidObj(v T, printErrors bool) bool {
	// st.tst.Helper()
	valJson, err := util.Jsonify(v)
	if err != nil {
		panic(err)
	}
	return st.AssertValid(valJson, printErrors)
}

func (st *SchemaTester[T]) AssertValid(val string, printErrors bool) bool {
	valLoader := gojsonschema.NewStringLoader(val)
	result, err := st.schema.Validate(valLoader)
	if err != nil {
		panic(err)
	}

	if result.Valid() {
		return true
	}

	if printErrors {
		st.tst.Errorf("SchemaFail: <%s>", val)
		for _, err := range result.Errors() {
			// Err implements the ResultError interface
			st.tst.Errorf("- %s", err)
		}
	}

	return false
}

func (st *SchemaTester[T]) AssertMatchExpected(actual T, expected T, printComplaint bool) bool {
	actualJson, err := util.Jsonify(actual)
	if err != nil {
		panic(err)
	}

	expJson, err := util.Jsonify(expected)
	if err != nil {
		panic(err)
	}

	if expJson != actualJson {

		if printComplaint {
			st.tst.Errorf("Mismatch. Expected=%s Got=%s", expJson, actualJson)
		}

		return false
	}

	return true
}

func (st *SchemaTester[T]) Parse(str string) (T, error) {

	var obj T

	if err := yaml.UnmarshalStrict([]byte(str), &obj); err != nil {
		return *new(T), err
	}

	return obj, nil
}
