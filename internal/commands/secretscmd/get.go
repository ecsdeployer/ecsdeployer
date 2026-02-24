package secretscmd

import (
	"fmt"
	"path/filepath"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/util/cmdutil"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/cobra"
)

type getCmdRunner struct {
	configFile   string
	outputFormat string
	valueOnly    bool
}

func newGetCmd() *cobra.Command {

	runner := &getCmdRunner{}
	cmd := &cobra.Command{
		Use:     "get VARIABLE_NAME [VARIABLE_NAME...]",
		Short:   "Get the value of a specific secret",
		PreRunE: runner.PreRunE,
		RunE:    runner.RunE,
		Args:    cobra.MatchAll(cobra.MinimumNArgs(1), cobra.MaximumNArgs(10)),
	}

	cmdutil.FlagConfigFile(cmd, &runner.configFile)

	cmd.Flags().StringVarP(&runner.outputFormat, "format", "f", outputFormatDotEnv, "The output format to use (dotenv, plain)")
	cmd.Flags().BoolVar(&runner.valueOnly, "bare", false, "Only return the value. You can only provide a single VARIABLE_NAME when using this.")

	return cmd
}

func (r *getCmdRunner) PreRunE(cmd *cobra.Command, args []string) error {
	if r.valueOnly && len(args) > 1 {
		return cmdutil.NewUserErrorf(`The --bare flag can only be used with a single VARIABLE.`)
	}
	return nil
}

func (r *getCmdRunner) RunE(cmd *cobra.Command, args []string) error {

	_, ssmPrefix, err := loadProject(cmd.Context(), r.configFile)
	if err != nil {
		return err
	}

	ssmClient := awsclients.SSMClient()

	names := make([]string, 0, len(args))
	for _, v := range args {
		names = append(names, ssmPrefix+v)
	}

	resp, err := ssmClient.GetParameters(cmd.Context(), &ssm.GetParametersInput{
		Names:          names,
		WithDecryption: new(true),
	})
	if err != nil {
		return err
	}

	for _, parameter := range resp.Parameters {
		// just want the last part of the name
		name := filepath.Base(*parameter.Name)

		if r.valueOnly {
			fmt.Fprint(cmd.OutOrStdout(), *parameter.Value)
			continue
		}

		switch r.outputFormat {
		case outputFormatPlain:
			fmt.Fprintf(cmd.OutOrStdout(), "%s: %q\n", name, *parameter.Value)

		default:
			// dotenv
			fmt.Fprintf(cmd.OutOrStdout(), "%s=%q\n", name, *parameter.Value)
		}
	}

	return nil
}
