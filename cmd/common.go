package cmd

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/spf13/cobra"
)

const (
	paramConfigFile = "config"
	paramAppVersion = "app-version"
	paramImageTag   = "tag"
	paramImage      = "image"
	// paramImageTag   = "image-tag"
)

type commonOpts struct {
	config     string
	appVersion string
	imageTag   string
	imageUri   string
}

func setCommonFlags(cmd *cobra.Command, common *commonOpts) {
	cmd.Flags().StringVarP(&common.config, paramConfigFile, "c", "", "Configuration file to check")
	cmd.Flags().StringVar(&common.appVersion, paramAppVersion, "", "Set the application version. Useful for templates")
	cmd.Flags().StringVar(&common.imageTag, paramImageTag, "", "Specify a custom image tag to use.")
	cmd.Flags().StringVar(&common.imageUri, paramImage, "", "Specify a container image URI.")

	_ = cmd.Flags().SetAnnotation(paramConfigFile, cobra.BashCompFilenameExt, []string{"yaml", "yml"})
}

func setupContextCommon(ctx *config.Context, options commonOpts) {
	ctx.Version = options.appVersion
	ctx.ImageTag = options.imageTag
	if options.imageUri != "" {
		ctx.ImageUriRef = options.imageUri
	}

	if ctx.Project.StageName != nil {
		ctx.Stage = *ctx.Project.StageName
	}
}
