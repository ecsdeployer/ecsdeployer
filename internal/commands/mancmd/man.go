package mancmd

import (
	"fmt"

	mcoral "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
)

func New() *cobra.Command {

	return &cobra.Command{
		Use:                   "man",
		Short:                 "Generates ECS Deployer's command line manpages",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Hidden:                true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			manPage, err := mcoral.NewManPage(1, cmd.Root())
			if err != nil {
				return err
			}

			// _, err = fmt.Fprint(os.Stdout, manPage.Build(roff.NewDocument()))
			_, err = fmt.Fprint(cmd.OutOrStdout(), manPage.Build(roff.NewDocument()))
			return err
		},
	}
}
