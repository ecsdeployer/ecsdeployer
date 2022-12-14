package cmd

import (
	"errors"
	"fmt"
	"time"

	"github.com/caarlos0/log"
	"github.com/spf13/cobra"
	cobracompletefig "github.com/withfig/autocomplete-tools/integrations/cobra"
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
	exit  func(int)
}

type cmdMetadata struct {
	version string
}

func newRootCmd(version string, exit func(int)) *rootCmd {
	root := &rootCmd{
		exit: exit,
	}

	metadata := &cmdMetadata{
		version: version,
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
			if root.debug {
				log.SetLevel(log.DebugLevel)
				log.Debug("debug logs enabled")
			}
		},
	}

	cmd.PersistentFlags().BoolVar(&root.debug, "debug", false, "Enable debug mode")

	cmd.AddCommand(
		newDeployCmd(metadata).cmd,
		newCheckCmd(metadata).cmd,
		newSchemaCmd(metadata).cmd,
		newManCmd(metadata).cmd,
		newDocsCmd(metadata).cmd,
		newCleanCmd(metadata).cmd,
		newInfoCmd(metadata).cmd,
		cobracompletefig.CreateCompletionSpecCommand(cobracompletefig.Opts{Visible: false}),
	)

	root.cmd = cmd
	return root
}

func timedRunE(verb string, runef func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		start := time.Now()

		log.Infof(fmt.Sprintf("starting %s...", verb))

		if err := runef(cmd, args); err != nil {
			return wrapError(err, fmt.Sprintf("%s failed after %s", verb, time.Since(start).Truncate(time.Second)))
		}

		log.Infof(fmt.Sprintf("%s succeeded after %s", verb, time.Since(start).Truncate(time.Second)))
		return nil
	}
}
