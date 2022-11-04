package cmd

import (
	"io"
	"time"

	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

type deployCmd struct {
	cmd  *cobra.Command
	opts deployOpts
}

type deployOpts struct {
	config   string
	quiet    bool
	timeout  time.Duration
	version  string
	imageTag string
	imageUri string
	metadata *cmdMetadata
	// noValidate bool
}

func newDeployCmd(metadata *cmdMetadata) *deployCmd {
	root := &deployCmd{}
	root.opts.metadata = metadata
	cmd := &cobra.Command{
		Use:           "deploy",
		Short:         "Deploys application",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE: timedRunE("deploy", func(cmd *cobra.Command, args []string) error {
			if root.opts.quiet {
				log.Log = log.New(io.Discard)
			}

			opts := &configLoaderExtras{
				configFile:  root.opts.config,
				appVersion:  root.opts.version,
				imageTag:    root.opts.imageTag,
				imageUri:    root.opts.imageUri,
				timeout:     root.opts.timeout,
				cmdMetadata: root.opts.metadata,
			}

			err := stepRunner(opts, stepRunModeDeploy)
			if err != nil {
				return err
			}
			return nil
		}),
	}

	cmd.Flags().BoolVarP(&root.opts.quiet, "quiet", "q", false, "Quiet mode: no output")
	cmd.Flags().DurationVar(&root.opts.timeout, "timeout", 2*time.Hour, "Timeout for the entire deploy process")

	setCommonFlags(cmd, &root.opts.config, &root.opts.version, &root.opts.imageTag, &root.opts.imageUri)
	// cmd.Flags().BoolVar(&root.opts.noValidate, "no-validate", false, "Skips validating the config file against the schema")
	// _ = cmd.Flags().SetAnnotation("config", cobra.BashCompFilenameExt, []string{"yaml", "yml"})

	root.cmd = cmd
	return root
}
