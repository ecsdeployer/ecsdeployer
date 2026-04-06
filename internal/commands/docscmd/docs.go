package docscmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func New() *cobra.Command {
	return &cobra.Command{
		Use:                   "docs",
		Short:                 "Generates ECS Deployer's command line docs",
		SilenceUsage:          true,
		DisableFlagsInUseLine: true,
		Hidden:                true,
		Args:                  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Root().DisableAutoGenTag = true
			return doc.GenMarkdownTreeCustom(cmd.Root(), "www/docs/cmd", func(_ string) string {
				return ""
			}, func(s string) string {
				return s
			})
		},
	}
}
