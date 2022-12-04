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
		cmd := newSchemaCmd(defaultCmdMetadata()).cmd
		dir := t.TempDir()
		destination := path.Join(dir, "schema.json")
		cmd.SetArgs([]string{"--output", destination})

		_, _, err := executeCmdAndReturnOutput(cmd)
		require.NoError(t, err)

		outFile, err := os.Open(destination)
		require.NoError(t, err)

		schema := map[string]interface{}{}
		require.NoError(t, json.NewDecoder(outFile).Decode(&schema))
		require.Equal(t, "https://json-schema.org/draft/2020-12/schema", schema["$schema"].(string))
	})

	t.Run("outputs to stdout", func(t *testing.T) {
		cmd := newSchemaCmd(defaultCmdMetadata()).cmd
		cmd.SetArgs([]string{"--output", "-"})
		// require.NoError(t, cmd.Execute())

		osOut, _, err := executeCmdAndReturnOutput(cmd)
		require.NoError(t, err)

		schema := map[string]interface{}{}
		require.NoError(t, json.NewDecoder(strings.NewReader(osOut)).Decode(&schema))
		require.Equal(t, "https://json-schema.org/draft/2020-12/schema", schema["$schema"].(string))
	})
}
