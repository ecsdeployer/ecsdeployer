package cmd

import "github.com/spf13/cobra"

const (
	paramConfigFile = "config"
	paramAppVersion = "app-version"
	paramImageTag   = "tag"
	paramImage      = "image"
	// paramImageTag   = "image-tag"
)

func setCommonFlags(cmd *cobra.Command, configPtr, verPtr, imgTagPtr, imgPtr *string) {
	cmd.Flags().StringVarP(configPtr, paramConfigFile, "c", "", "Configuration file to check")
	cmd.Flags().StringVar(verPtr, paramAppVersion, "", "Set the application version. Useful for templates")
	cmd.Flags().StringVar(imgTagPtr, paramImageTag, "", "Specify a custom image tag to use.")
	cmd.Flags().StringVar(imgPtr, paramImage, "", "Specify a container image URI.")

	_ = cmd.Flags().SetAnnotation(paramConfigFile, cobra.BashCompFilenameExt, []string{"yaml", "yml"})
}
