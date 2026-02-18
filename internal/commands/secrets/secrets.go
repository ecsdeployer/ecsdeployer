package secrets

import "github.com/spf13/cobra"

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secrets",
		Short: "Manage application secrets in SSM",
	}

	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newSetCmd())
	cmd.AddCommand(newGetCmd())

	return cmd
}
