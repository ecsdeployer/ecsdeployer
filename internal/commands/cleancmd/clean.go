package cleancmd

import (
	"io"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/middleware/errhandler"
	"ecsdeployer.com/ecsdeployer/internal/middleware/logging"
	"ecsdeployer.com/ecsdeployer/internal/middleware/skip"
	"ecsdeployer.com/ecsdeployer/internal/pipeline"
	"ecsdeployer.com/ecsdeployer/internal/util/cmdutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/caarlos0/ctrlc"
	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
)

type cleanRunner struct {
	cmdutil.CommonOptions
	quiet   bool
	timeout time.Duration
}

func New() *cobra.Command {
	runner := &cleanRunner{}
	cmd := &cobra.Command{
		Use:   "clean",
		Short: "Runs the cleanup step only. Skips actual deployment",
		Long: `Use this command to purge any unused services, cronjobs, task definitions, etc
from your environment that are no longer being referenced in your configuration file.
`,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE:          cmdutil.TimedRunE("clean", runner.RunE),
	}

	cmd.Flags().BoolVarP(&runner.quiet, "quiet", "q", false, "Quiet mode: no output")
	cmd.Flags().DurationVar(&runner.timeout, "timeout", 30*time.Minute, "Timeout for the entire cleanup process")

	cmdutil.SetCommonFlags(cmd, &runner.CommonOptions)

	return cmd
}

func (r *cleanRunner) RunE(cmd *cobra.Command, args []string) error {
	if r.quiet {
		log.Log = log.New(io.Discard)
	}

	ctx, err := r.cleanProject()
	if err != nil {
		return err
	}
	cmdutil.DeprecateWarn(ctx)
	return nil
}

func (r *cleanRunner) cleanProject() (*config.Context, error) {
	cfg, err := cmdutil.LoadConfig(r.ConfigFile)
	if err != nil {
		return nil, err
	}
	ctx, cancel := config.NewWithTimeout(cfg, r.timeout)
	defer cancel()
	cmdutil.SetCommonContext(ctx, r.CommonOptions)
	ctx.CleanOnlyFlow = true
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
