package configschema_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/require"
)

func TestExports(t *testing.T) {
	// require.Equal(t, "#/$defs/StringLike", configschema.StringLikeRef)
	// require.Equal(t, "#/$defs/StringLikeWithBlank", configschema.StringLikeWithBlankRef)

	require.IsType(t, &jsonschema.Schema{}, configschema.StringLike)
	require.IsType(t, &jsonschema.Schema{}, configschema.StringLikeWithBlank)
}

func TestSchemaNamer(t *testing.T) {

	type bigHolder struct {
		TaskDefaults   *config.FargateDefaults
		Network        *config.NetworkConfiguration
		Storage        *config.StorageSpec
		CPUShares      *config.CpuSpec
		Memory         *config.MemorySpec
		RoleRef        *config.RoleArn
		ClusterRef     *config.ClusterArn
		TargetGroupRef *config.TargetGroupArn
		Service        *config.Service
	}

	tables := []struct {
		expectedName string
		intf         interface{}
	}{
		// different name
		{"TaskDefaults", config.FargateDefaults{}},
		{"Network", config.NetworkConfiguration{}},
		{"Storage", config.StorageSpec(1)},
		{"CPUShares", config.CpuSpec(1)},
		{"Memory", config.MemorySpec{}},
		{"RoleRef", config.RoleArn{}},
		{"ClusterRef", config.ClusterArn{}},
		{"TargetGroupRef", config.TargetGroupArn{}},

		// Same name
		{"Service", config.Service{}},
	}

	schema := configschema.GenerateSchema(&bigHolder{})

	for _, table := range tables {
		t.Run(table.expectedName, func(t *testing.T) {

			prop, ok := schema.Properties.Get(table.expectedName)
			require.True(t, ok)
			property := prop.(*jsonschema.Schema)
			require.Equal(t, fmt.Sprintf("#/$defs/%s", table.expectedName), property.Ref)

			expectedSchema := configschema.GenerateSchema(table.intf)

			expectedSchemaJson := jsonSchemaWithoutSpecials(expectedSchema)
			defSchemaJson := jsonSchemaWithoutSpecials(schema.Definitions[table.expectedName])

			require.JSONEq(t, expectedSchemaJson, defSchemaJson)

		})
	}
}

func TestGenerateSchema(t *testing.T) {
	tables := []struct {
		entity interface{}
	}{
		{&config.Project{}},
		{&config.Service{}},
		{&config.PreDeployTask{}},
		{&config.ConsoleTask{}},
		{&config.CronJob{}},
		{&config.SSMImport{}},
		{&config.Settings{}},
		{&config.ClusterArn{}},
		{&config.RoleArn{}},
		{&config.NameTemplates{}},
	}

	for _, table := range tables {
		schema := configschema.GenerateSchema(table.entity)
		require.NotNil(t, schema)
	}
}

// strips the $xxxx fields from the schema for easier comparisons
func jsonSchemaWithoutSpecials(schema *jsonschema.Schema) string {
	initialJson, _ := util.Jsonify(schema)

	temp := make(map[string]interface{})
	if err := json.Unmarshal([]byte(initialJson), &temp); err != nil {
		panic(err)
	}

	for k := range temp {
		if strings.HasPrefix(k, "$") {
			delete(temp, k)
		}
	}

	cleanJson, _ := util.Jsonify(temp)

	return cleanJson
}
