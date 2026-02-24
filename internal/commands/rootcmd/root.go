package rootcmd

import (
	"ecsdeployer.com/ecsdeployer/internal/commands/checkcmd"
	"ecsdeployer.com/ecsdeployer/internal/commands/cleancmd"
	"ecsdeployer.com/ecsdeployer/internal/commands/deploycmd"
	"ecsdeployer.com/ecsdeployer/internal/commands/docscmd"
	"ecsdeployer.com/ecsdeployer/internal/commands/infocmd"
	"ecsdeployer.com/ecsdeployer/internal/commands/mancmd"
	"ecsdeployer.com/ecsdeployer/internal/commands/schemacmd"
	"ecsdeployer.com/ecsdeployer/internal/commands/secretscmd"
	"ecsdeployer.com/ecsdeployer/internal/util/cmdutil"
	"ecsdeployer.com/ecsdeployer/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/webdestroya/go-log"
	cobracompletefig "github.com/withfig/autocomplete-tools/integrations/cobra"
)

type rootCmd struct {
	debug bool
	trace bool
}

func New() *cobra.Command {

	cobra.EnableTraverseRunHooks = true

	runner := &rootCmd{}

	cmd := &cobra.Command{
		Use:   "ecsdeployer",
		Short: "Deploy applications to Fargate",
		Long: `ECS Deployer allows you to easily deploy containerized applications to AWS ECS Fargate.
It simplifies the process of creating task definitions, running pre-deployment tasks, setting up scheduled jobs,
as well as service orchestration.

Applications can easily and securely be deployed with a simple GitHub Action.

Check out our website for more information, examples and documentation: https://ecsdeployer.com/
`,
		Version:           version.String(),
		SilenceUsage:      true,
		SilenceErrors:     true,
		Args:              cobra.NoArgs,
		PersistentPreRunE: runner.PersistentPreRunE,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}

	cmd.PersistentFlags().BoolVar(&runner.debug, "debug", false, "Enable debug mode")
	cmd.PersistentFlags().BoolVar(&runner.trace, "trace", false, "Enable trace mode")
	_ = cmd.PersistentFlags().MarkHidden("trace")

	cmd.SetFlagErrorFunc(rootFlagErrorFunc)

	cmd.PersistentFlags().BoolP("help", "h", false, "Print usage")
	cmd.PersistentFlags().Lookup("help").Hidden = true

	cmd.AddCommand(
		deploycmd.New(), // newDeployCmd().cmd,
		checkcmd.New(),  // newCheckCmd().cmd,
		schemacmd.New(), // newSchemaCmd().cmd,
		mancmd.New(),    // newManCmd().cmd,
		docscmd.New(),   // newDocsCmd().cmd,
		cleancmd.New(),  // newCleanCmd().cmd,
		infocmd.New(),   // newInfoCmd().cmd,
		secretscmd.New(),
		cobracompletefig.CreateCompletionSpecCommand(cobracompletefig.Opts{Visible: false}),
	)

	return cmd
}

func (r *rootCmd) PersistentPreRunE(cmd *cobra.Command, args []string) error {
	log.Strings[log.DebugLevel] = "%"
	if r.trace {
		log.SetLevel(log.TraceLevel)
		log.Debug("trace logs enabled")
	} else if r.debug {
		log.SetLevel(log.DebugLevel)
		log.Debug("debug logs enabled")
	}
	return nil
}

func rootFlagErrorFunc(cmd *cobra.Command, err error) error {
	if err == pflag.ErrHelp {
		return err
	}
	return cmdutil.FlagErrorWrap(err)
}
