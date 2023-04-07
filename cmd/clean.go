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
	"github.com/spf13/cobra"
	"github.com/webdestroya/go-log"
)

type cleanCmd struct {
	cmd  *cobra.Command
	opts cleanOpts
}

type cleanOpts struct {
	commonOpts
	quiet   bool
	timeout time.Duration
}

func newCleanCmd() *cleanCmd {
	root := &cleanCmd{}
	// root.opts.metadata = metadata
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Runs the cleanup step only. Skips actual deployment",
		Long: `Use this command to purge any unused services, cronjobs, task definitions, etc 
from your environment that are no longer being referenced in your configuration file.
`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE: timedRunE("clean", func(cmd *cobra.Command, args []string) error {
			if root.opts.quiet {
				log.Log = log.New(io.Discard)
			}
			ctx, err := cleanProject(root.opts)
			if err != nil {
				return err
			}
			deprecateWarn(ctx)
			return nil
		}),
	}

	cmd.Flags().BoolVarP(&root.opts.quiet, "quiet", "q", false, "Quiet mode: no output")
	cmd.Flags().DurationVar(&root.opts.timeout, "timeout", 30*time.Minute, "Timeout for the entire cleanup process")

	setCommonFlags(cmd, &root.opts.commonOpts)

	root.cmd = cmd
	return root
}

func cleanProject(options cleanOpts) (*config.Context, error) {
	cfg, err := loadConfig(options.config)
	if err != nil {
		return nil, err
	}
	ctx, cancel := config.NewWithTimeout(cfg, options.timeout)
	defer cancel()
	setupCleanContext(ctx, options)
	return ctx, ctrlc.Default.Run(ctx, func() error {
		for _, step := range pipeline.CleanupPipeline {
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

func setupCleanContext(ctx *config.Context, options cleanOpts) {
	setupContextCommon(ctx, options.commonOpts)
	ctx.CleanOnlyFlow = true
}
