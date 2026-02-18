package secrets

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/internal/util/cmdutil"
	"github.com/spf13/cobra"
)

type setCmdRunner struct {
	configFile string
	unset      bool
	force      bool
}

func newSetCmd() *cobra.Command {

	runner := &setCmdRunner{}

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set (or delete) a secret",
		RunE:  runner.RunE,
	}

	cmdutil.FlagConfigFile(cmd, &runner.configFile)
	cmd.Flags().BoolVar(&runner.unset, "unset", false, "Removes a variable entirely")
	cmd.Flags().BoolVar(&runner.force, "force", false, "Do not ask for confirmation")

	return cmd
}

func (r *setCmdRunner) RunE(cmd *cobra.Command, args []string) error {
	return errors.New("NOT FINISHED")
}
