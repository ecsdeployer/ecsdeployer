package cmdutil

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/spf13/cobra"
)

type CommonOptions struct {
	ConfigFile string
	AppVersion string
	ImageTag   string
	ImageUri   string
}

func SetCommonFlags(cmd *cobra.Command, common *CommonOptions) {
	FlagConfigFile(cmd, &common.ConfigFile)
	cmd.Flags().StringVar(&common.AppVersion, paramAppVersion, "", "Set the application version. Useful for templates")
	cmd.Flags().StringVar(&common.ImageTag, paramImageTag, "", "Specify a custom image tag to use.")
	cmd.Flags().StringVar(&common.ImageUri, paramImage, "", "Specify a container image URI.")
}

func SetCommonContext(ctx *config.Context, common CommonOptions) {
	ctx.Version = common.AppVersion
	ctx.ImageTag = common.ImageTag
	if common.ImageUri != "" {
		ctx.ImageUriRef = common.ImageUri
	}

	if ctx.Project.StageName != nil {
		ctx.Stage = *ctx.Project.StageName
	}
}
