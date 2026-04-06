package secretscmd

import (
	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/usererr"
	"ecsdeployer.com/ecsdeployer/internal/util/cmdutil"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/spf13/cobra"
	"github.com/webdestroya/go-log"
)

type deleteCmdRunner struct {
	configFile string
	varName    string
}

func newDeleteCmd() *cobra.Command {

	runner := &deleteCmdRunner{}

	cmd := &cobra.Command{
		Use:     "delete VARIABLE_NAME",
		Short:   "Delete a secret",
		RunE:    runner.RunE,
		PreRunE: runner.PreRunE,
		Args:    cobra.ExactArgs(1),
	}

	cmdutil.FlagConfigFile(cmd, &runner.configFile)

	return cmd
}
func (r *deleteCmdRunner) PreRunE(cmd *cobra.Command, args []string) error {

	r.varName = args[0]

	return nil
}

func (r *deleteCmdRunner) RunE(cmd *cobra.Command, args []string) error {

	_, ssmPrefix, err := loadProject(cmd.Context(), r.configFile)
	if err != nil {
		return err
	}

	ssmClient := awsclients.SSMClient()

	paramName := ssmPrefix + r.varName

	if _, err = ssmClient.DeleteParameter(cmd.Context(), &ssm.DeleteParameterInput{
		Name: &paramName,
	}); err != nil {
		// if _, ok := errors.AsType[*ssmTypes.ParameterNotFound](err); ok {

		// 	log.WithField("param", paramName).Error("parameter not found")
		// 	return nil
		// }

		return usererr.Wrap(err)
	}
	log.WithField("param", paramName).Info("parameter deleted")

	log.Info("You must redeploy the project so containers pick up the new value.")

	return nil
}
