package cmd

import (
	"io"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/middleware/errhandler"
	"ecsdeployer.com/ecsdeployer/internal/middleware/logging"
	"ecsdeployer.com/ecsdeployer/internal/middleware/skip"
	"ecsdeployer.com/ecsdeployer/internal/pipeline"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/caarlos0/ctrlc"
	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

type deployCmd struct {
	cmd  *cobra.Command
	opts deployOpts
}

type deployOpts struct {
	commonOpts
	quiet   bool
	timeout time.Duration
}

func newDeployCmd() *deployCmd {
	root := &deployCmd{}
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

			ctx, err := deployProject(root.opts)
			if err != nil {
				return err
			}
			deprecateWarn(ctx)
			return nil
		}),
	}

	cmd.Flags().BoolVarP(&root.opts.quiet, "quiet", "q", false, "Quiet mode: no output")
	cmd.Flags().DurationVar(&root.opts.timeout, "timeout", 2*time.Hour, "Timeout for the entire deploy process")

	setCommonFlags(cmd, &root.opts.commonOpts)
	// cmd.Flags().BoolVar(&root.opts.noValidate, "no-validate", false, "Skips validating the config file against the schema")
	// _ = cmd.Flags().SetAnnotation("config", cobra.BashCompFilenameExt, []string{"yaml", "yml"})

	root.cmd = cmd
	return root
}

func deployProject(options deployOpts) (*config.Context, error) {
	cfg, err := loadConfig(options.config)
	if err != nil {
		return nil, err
	}
	ctx, cancel := config.NewWithTimeout(cfg, options.timeout)
	defer cancel()
	setupDeployContext(ctx, options)
	return ctx, ctrlc.Default.Run(ctx, func() error {
		for _, step := range pipeline.DeploymentPipeline {
			if err := skip.Maybe(
				step,
				logging.Log(
					step.String(),
					errhandler.Handle(step.Run),
				),
			)(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func setupDeployContext(ctx *config.Context, options deployOpts) {
	setupContextCommon(ctx, options.commonOpts)
}
