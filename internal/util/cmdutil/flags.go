package cmdutil

import "github.com/spf13/cobra"

const (
	paramConfigFile = "config"
	paramAppVersion = "app-version"
	paramImageTag   = "tag"
	paramImage      = "image"
)

func FlagConfigFile(cmd *cobra.Command, dest *string) {
	cmd.Flags().StringVarP(dest, paramConfigFile, "c", "", "Configuration `file` to check")

	_ = cmd.Flags().SetAnnotation(paramConfigFile, cobra.BashCompFilenameExt, []string{"yaml", "yml"})
}
