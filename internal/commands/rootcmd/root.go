package rootcmd

import (
	"ecsdeployer.com/ecsdeployer/internal/util/cmdutil"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

func New() *cobra.Command {

	cmd := &cobra.Command{}

	return cmd
}

func rootFlagErrorFunc(cmd *cobra.Command, err error) error {
	if err == pflag.ErrHelp {
		return err
	}
	return cmdutil.FlagErrorWrap(err)
}
