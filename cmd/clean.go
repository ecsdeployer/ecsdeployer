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
	log "github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

type cleanCmd struct {
	cmd  *cobra.Command
	opts cleanOpts
}

type cleanOpts struct {
	config   string
	quiet    bool
	version  string
	imageUri string
	imageTag string
	metadata *cmdMetadata
}

func newCleanCmd(metadata *cmdMetadata) *cleanCmd {
	root := &cleanCmd{}
	root.opts.metadata = metadata
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
			opts := &configLoaderExtras{
				configFile:  root.opts.config,
				appVersion:  root.opts.version,
				imageTag:    root.opts.imageTag,
				imageUri:    root.opts.imageUri,
				cmdMetadata: root.opts.metadata,
			}

			err := stepRunner(opts, stepRunModeCleanup)
			if err != nil {
				return err
			}
			return nil
		}),
	}

	cmd.Flags().BoolVarP(&root.opts.quiet, "quiet", "q", false, "Quiet mode: no output")

	setCommonFlags(cmd, &root.opts.config, &root.opts.version, &root.opts.imageTag, &root.opts.imageUri)

	root.cmd = cmd
	return root
}

func cleanProject(options cleanOpts) (*config.Context, error) {
	cfg, err := loadConfig(options.config)
	if err != nil {
		return nil, err
	}
	ctx, cancel := config.NewWithTimeout(cfg, 30*time.Minute)
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
}
