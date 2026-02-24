package secretscmd

import (
	"fmt"
	"path/filepath"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/util/cmdutil"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/cobra"
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
	ctx, ssmPrefix, err := loadProject(cmd.Context(), r.configFile)
	if err != nil {
		return err
	}

	ssmClient := awsclients.SSMClient()

	request := &ssm.GetParametersByPathInput{
		Path:           &ssmPrefix,
		WithDecryption: new(true),
		Recursive:      ctx.Project.Settings.SSMImport.Recursive,
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
