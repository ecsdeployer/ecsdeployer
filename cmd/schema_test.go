package cmd

// these exist here to prevent import cycle

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestSchemaCmd(t *testing.T) {
	cmd := newSchemaCmd(defaultCmdMetadata()).cmd
	dir := t.TempDir()
	destination := path.Join(dir, "schema.json")
	cmd.SetArgs([]string{"--output", destination})
	require.NoError(t, cmd.Execute())

	outFile, err := os.Open(destination)
	require.NoError(t, err)

	schema := map[string]interface{}{}
	require.NoError(t, json.NewDecoder(outFile).Decode(&schema))
	require.Equal(t, "https://json-schema.org/draft/2020-12/schema", schema["$schema"].(string))
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

		if schema == nil {
			t.Error("expected schema object to not be nil")
			continue
		}

		_, err := util.Jsonify(schema)
		if schema == nil {
			t.Errorf("expected schema to be able to jsonify: %s", err)
			continue
		}
	}
}
