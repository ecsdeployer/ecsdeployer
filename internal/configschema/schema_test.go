package configschema_test

import (
	"reflect"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/invopop/jsonschema"
	"github.com/stretchr/testify/require"
)

func TestExports(t *testing.T) {
	require.Equal(t, "#/$defs/StringLike", configschema.StringLikeRef)
	require.Equal(t, "#/$defs/StringLikeWithBlank", configschema.StringLikeWithBlankRef)

	require.IsType(t, &jsonschema.Schema{}, configschema.StringLike)
	require.IsType(t, &jsonschema.Schema{}, configschema.StringLikeWithBlank)
}

func TestSchemaNamer(t *testing.T) {
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

	for _, table := range tables {
		t.Run(table.expectedName, func(t *testing.T) {
			actualName := configschema.SchemaNamer(reflect.ValueOf(table.intf).Type())

			require.Equal(t, table.expectedName, actualName)

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
