package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/spf13/cobra"
)

type schemaCmd struct {
	cmd    *cobra.Command
	output string
}

func newSchemaCmd(metadata *cmdMetadata) *schemaCmd {
	root := &schemaCmd{}
	cmd := &cobra.Command{
		Use:           "jsonschema",
		Aliases:       []string{"schema"},
		Short:         "outputs ECS Deployer's JSON schema",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			schema := configschema.GenerateSchema(&config.Project{})
			bts, err := json.MarshalIndent(schema, " ", " ")
			if err != nil {
				return fmt.Errorf("failed to create jsonschema: %w", err)
			}
			if root.output == "-" {
				cmd.Println(string(bts))
				return nil
			}
			if err := os.MkdirAll(filepath.Dir(root.output), 0o755); err != nil {
				return fmt.Errorf("failed to write jsonschema file: %w", err)
			}
			if err := os.WriteFile(root.output, bts, 0o600); err != nil {
				return fmt.Errorf("failed to write jsonschema file: %w", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&root.output, "output", "o", "-", "where to save the json schema")

	root.cmd = cmd
	return root
}
