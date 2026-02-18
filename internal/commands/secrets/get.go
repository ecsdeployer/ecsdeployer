package secrets

import (
	"ecsdeployer.com/ecsdeployer/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

type getCmdRunner struct {
	configFile string
}

func newGetCmd() *cobra.Command {

	runner := &getCmdRunner{}
	cmd := &cobra.Command{
		Use:   "get NAME",
		Short: "Get the value of a specific secret",
		RunE:  runner.RunE,
		Args:  cobra.ExactArgs(1),
	}

	cmdutil.FlagConfigFile(cmd, &runner.configFile)

	return cmd
}

func (r *getCmdRunner) RunE(cmd *cobra.Command, args []string) error {
	return nil
}
