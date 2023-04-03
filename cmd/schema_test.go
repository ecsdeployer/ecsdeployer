package cmd

// these exist here to prevent import cycle

import (
	"encoding/json"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSchemaCmd(t *testing.T) {

	t.Run("writes file", func(t *testing.T) {
		if testing.Short() {
			t.SkipNow()
			return
		}
		dir := t.TempDir()
		destination := path.Join(dir, "schema.json")

		result := runCommand("schema", "--output", destination)
		require.NoError(t, result.err)

		outFile, err := os.Open(destination)
		require.NoError(t, err)

		schema := map[string]interface{}{}
		require.NoError(t, json.NewDecoder(outFile).Decode(&schema))
		require.Equal(t, "https://json-schema.org/draft/2020-12/schema", schema["$schema"].(string))
	})

	t.Run("outputs to stdout", func(t *testing.T) {
		result := runCommand("schema", "--output", "-")
		require.NoError(t, result.err)

		schema := map[string]interface{}{}
		require.NoError(t, json.NewDecoder(strings.NewReader(result.stdout)).Decode(&schema))
		require.Equal(t, "https://json-schema.org/draft/2020-12/schema", schema["$schema"].(string))
	})
}
