package cmd

import (
	"errors"
	"fmt"
	"time"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	"github.com/webdestroya/go-log"
	cobracompletefig "github.com/withfig/autocomplete-tools/integrations/cobra"
)

var (
	boldStyle = lipgloss.NewStyle().Bold(true)
)

func Execute(version string, exit func(int), args []string) {
	newRootCmd(version, exit).Execute(args)
}

func (cmd *rootCmd) Execute(args []string) {
	cmd.cmd.SetArgs(args)

	if err := cmd.cmd.Execute(); err != nil {
		code := 1
		msg := "command failed"
		eerr := &exitError{}
		if errors.As(err, &eerr) {
			code = eerr.code
			if eerr.details != "" {
				msg = eerr.details
			}
		}
		log.WithError(err).Error(msg)
		cmd.exit(code)
	}
}

type rootCmd struct {
	cmd   *cobra.Command
	debug bool
	trace bool
	exit  func(int)
}

func newRootCmd(version string, exit func(int)) *rootCmd {
	root := &rootCmd{
		exit: exit,
	}

	cmd := &cobra.Command{
		Use:   "ecsdeployer",
		Short: "Deploy applications to Fargate",
		Long: `ECS Deployer allows you to easily deploy containerized applications to AWS ECS Fargate.
It simplifies the process of creating task definitions, running pre-deployment tasks, setting up scheduled jobs,
as well as service orchestration.

Applications can easily and securely be deployed with a simple GitHub Action.

Check out our website for more information, examples and documentation: https://ecsdeployer.com/
`,
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
		Args:          cobra.NoArgs,
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {

			log.Strings[log.DebugLevel] = "%"
			if root.trace {
				log.SetLevel(log.TraceLevel)
				log.Debug("trace logs enabled")
			} else if root.debug {
				log.SetLevel(log.DebugLevel)
				log.Debug("debug logs enabled")
			}
		},
	}

	cmd.PersistentFlags().BoolVar(&root.debug, "debug", false, "Enable debug mode")
	cmd.PersistentFlags().BoolVar(&root.trace, "trace", false, "Enable trace mode")
	_ = cmd.PersistentFlags().MarkHidden("trace")

	cmd.AddCommand(
		newDeployCmd().cmd,
		newCheckCmd().cmd,
		newSchemaCmd().cmd,
		newManCmd().cmd,
		newDocsCmd().cmd,
		newCleanCmd().cmd,
		newInfoCmd().cmd,
		cobracompletefig.CreateCompletionSpecCommand(cobracompletefig.Opts{Visible: false}),
	)

	root.cmd = cmd
	return root
}

func deprecateWarn(ctx *config.Context) {
	if ctx.Deprecated {
		log.Warn(boldStyle.Render("you are using deprecated features, check the log above for information"))
	}
}

func timedRunE(verb string, runef func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		start := time.Now()

		log.Info(fmt.Sprintf("starting %s...", verb))

		if err := runef(cmd, args); err != nil {
			return wrapError(err, fmt.Sprintf("%s failed after %s", verb, time.Since(start).Truncate(time.Second)))
		}

		log.Info(fmt.Sprintf("%s succeeded after %s", verb, time.Since(start).Truncate(time.Second)))
		return nil
	}
}
