package schemacmd

import (
	"fmt"
	"os"
	"path/filepath"

	"ecsdeployer.com/ecsdeployer/internal/configschema"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/spf13/cobra"
)

type schemaCmd struct {
	output string
}

func New() *cobra.Command {
	runner := &schemaCmd{}
	cmd := &cobra.Command{
		Use:           "jsonschema",
		Aliases:       []string{"schema"},
		Short:         "outputs ECS Deployer's JSON schema",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE:          runner.RunE,
	}

	cmd.Flags().StringVarP(&runner.output, "output", "o", "-", "where to save the json schema")
	return cmd
}

func (r *schemaCmd) RunE(cmd *cobra.Command, args []string) error {
	schema := configschema.GenerateSchema(&config.Project{})
	bts, err := util.JsonifyPretty(schema)
	if err != nil {
		return fmt.Errorf("failed to create jsonschema: %w", err)
	}
	if r.output == "-" {
		fmt.Fprintln(cmd.OutOrStdout(), bts)
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(r.output), 0o755); err != nil {
		return fmt.Errorf("failed to write jsonschema file: %w", err)
	}
	if err := os.WriteFile(r.output, []byte(bts), 0o600); err != nil {
		return fmt.Errorf("failed to write jsonschema file: %w", err)
	}
	return nil
}
