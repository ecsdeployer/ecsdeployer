package secretscmd

import (
	"io"
	"os"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/util/cmdutil"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/spf13/cobra"
	"github.com/webdestroya/go-log"
)

type setCmdRunner struct {
	configFile string
	unset      bool
	force      bool
	valueFile  string
	useStdin   bool

	value   string
	varName string
	keyId   string
}

func newSetCmd() *cobra.Command {

	runner := &setCmdRunner{}

	cmd := &cobra.Command{
		Use:     "set VARIABLE_NAME {VALUE | --file filename | --stdin | --unset}",
		Short:   "Set (or delete) a secret",
		RunE:    runner.RunE,
		PreRunE: runner.PreRunE,
		Args:    cobra.RangeArgs(1, 2),
	}

	cmdutil.FlagConfigFile(cmd, &runner.configFile)
	cmd.Flags().BoolVar(&runner.unset, "unset", false, "Removes a variable entirely")
	cmd.Flags().BoolVar(&runner.force, "force", false, "Do not ask for confirmation")
	cmd.Flags().BoolVar(&runner.useStdin, "stdin", false, "Get value from stdin")
	cmd.Flags().StringVar(&runner.valueFile, "file", "", "Read value from `file`")

	cmd.Flags().StringVar(&runner.keyId, "keyid", "", "Specify the KMS `KeyID` to use. If not provided, will use the default.")

	return cmd
}
func (r *setCmdRunner) PreRunE(cmd *cobra.Command, args []string) error {

	optCount := 0
	if r.valueFile != "" {
		optCount++
	}
	if r.unset {
		optCount++
	}
	if len(args) > 1 {
		optCount++
	}
	if r.useStdin {
		optCount++
	}

	if r.valueFile == "-" {
		r.valueFile = ""
		r.useStdin = true
	}

	if optCount == 0 {
		return cmdutil.NewUserErrorf("You must provide either a value or the --unset flag")
	} else if optCount > 1 {
		return cmdutil.NewUserErrorf("You can only provide one of: VALUE, --file, --stdin, --unset")
	}

	r.varName = args[0]

	if len(args) > 1 {
		r.value = args[1]
	}

	if r.valueFile != "" {
		info, err := os.Stat(r.valueFile)
		if err != nil {
			return err
		}
		if info.Size() > 10*1024 {
			return cmdutil.NewUserErrorf("File provided exceeds 10KB, it will not fit into a Parameter")
		}

		val, err := os.ReadFile(r.valueFile)
		if err != nil {
			return err
		}
		r.value = string(val)
	}

	if r.useStdin {
		val, err := io.ReadAll(cmd.InOrStdin())
		if err != nil {
			return err
		}
		r.value = string(val)
	}

	return nil
}

func (r *setCmdRunner) RunE(cmd *cobra.Command, args []string) error {

	_, ssmPrefix, err := loadProject(cmd.Context(), r.configFile)
	if err != nil {
		return err
	}

	ssmClient := awsclients.SSMClient()

	paramName := ssmPrefix + r.varName

	if r.unset {

		if _, err = ssmClient.DeleteParameter(cmd.Context(), &ssm.DeleteParameterInput{
			Name: &paramName,
		}); err != nil {
			return err
		}
		log.WithField("param", paramName).Info("parameter deleted")

	} else {

		params := &ssm.PutParameterInput{
			Name:      &paramName,
			Value:     &r.value,
			Overwrite: new(true),
			Type:      ssmTypes.ParameterTypeSecureString,
		}
		if r.keyId != "" {
			params.KeyId = &r.keyId
		}

		if _, err := ssmClient.PutParameter(cmd.Context(), params); err != nil {
			return err
		}
		log.WithField("param", paramName).Info("parameter set")
	}

	log.Info("You must redeploy the project so containers pick up the new value.")

	return nil
}
