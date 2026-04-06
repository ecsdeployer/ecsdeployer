package deploycmd

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

type deployRunner struct {
	cmdutil.CommonOptions
	quiet   bool
	timeout time.Duration
}

func New() *cobra.Command {
	runner := &deployRunner{}
	cmd := &cobra.Command{
		Use:           "deploy",
		Short:         "Deploys application",
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		RunE:          cmdutil.TimedRunE("deploy", runner.RunE),
	}

	cmd.Flags().BoolVarP(&runner.quiet, "quiet", "q", false, "Quiet mode: no output")
	cmd.Flags().DurationVar(&runner.timeout, "timeout", 2*time.Hour, "Timeout for the entire deploy process")

	cmdutil.SetCommonFlags(cmd, &runner.CommonOptions)
	// cmd.Flags().BoolVar(&root.opts.noValidate, "no-validate", false, "Skips validating the config file against the schema")
	// _ = cmd.Flags().SetAnnotation("config", cobra.BashCompFilenameExt, []string{"yaml", "yml"})

	return cmd
}

func (r *deployRunner) RunE(cmd *cobra.Command, args []string) error {
	if r.quiet {
		log.Log = log.New(io.Discard)
	}

	ctx, err := r.deployProject()
	if err != nil {
		return err
	}
	cmdutil.DeprecateWarn(ctx)
	return nil
}

func (r *deployRunner) deployProject() (*config.Context, error) {
	cfg, err := cmdutil.LoadConfig(r.ConfigFile)
	if err != nil {
		return nil, err
	}
	ctx, cancel := config.NewWithTimeout(cfg, r.timeout)
	defer cancel()
	cmdutil.SetCommonContext(ctx, r.CommonOptions)
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
