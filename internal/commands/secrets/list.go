package secrets

import (
	"fmt"
	"path/filepath"
	"strings"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/internal/util/cmdutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/cobra"
)

const (
	outputFormatDotEnv = `dotenv`
	outputFormatPlain  = `plain`
)

type listCmdRunner struct {
	configFile string

	outputFormat string
}

func newListCmd() *cobra.Command {

	runner := &listCmdRunner{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List secrets in SSM that will be added to this project",
		RunE:  runner.RunE,
		Args:  cobra.NoArgs,
	}

	cmdutil.FlagConfigFile(cmd, &runner.configFile)

	cmd.Flags().StringVarP(&runner.outputFormat, "format", "f", outputFormatDotEnv, "The output format to use (dotenv, plain)")

	return cmd
}

func (r *listCmdRunner) RunE(cmd *cobra.Command, args []string) error {

	proj, err := cmdutil.LoadConfig(r.configFile)
	if err != nil {
		return err
	}

	if !proj.Settings.SSMImport.IsEnabled() {
		return fmt.Errorf(`SSM import is not enabled for this project, nothing to list.`)
	}

	ssmImport := *proj.Settings.SSMImport

	ctx := config.Wrap(cmd.Context(), proj)

	ssmPrefix, err := tmpl.New(ctx).Apply(ssmImport.GetPath())
	if err != nil {
		return err
	}

	// Trim any trailing slash, then add our own
	ssmPrefix = strings.TrimSuffix(ssmPrefix, "/") + "/"

	ssmClient := awsclients.SSMClient()

	request := &ssm.GetParametersByPathInput{
		Path:           &ssmPrefix,
		WithDecryption: new(true),
		Recursive:      ssmImport.Recursive,
	}

	paginator := ssm.NewGetParametersByPathPaginator(ssmClient, request, func(o *ssm.GetParametersByPathPaginatorOptions) {
		o.Limit = 10
	})

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx.Context)
		if err != nil {
			return fmt.Errorf("failed to load secrets: %w", err)
		}
		for _, parameter := range output.Parameters {
			// just want the last part of the name
			name := filepath.Base(*parameter.Name)

			switch r.outputFormat {
			case outputFormatPlain:
				fmt.Fprintf(cmd.OutOrStdout(), "%s: %q\n", name, *parameter.Value)

			default:
				// dotenv
				fmt.Fprintf(cmd.OutOrStdout(), "%s=%q\n", name, *parameter.Value)
			}

		}
	}

	return nil
}
