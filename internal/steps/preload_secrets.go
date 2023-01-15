package steps

import (
	"path/filepath"
	"strings"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

func PreloadSecretsStep(project *config.Project) *Step {
	return NewStep(&Step{
		Label:  "PreloadSecrets",
		Create: stepPreloadSecretsCreate,
	})
}

func stepPreloadSecretsCreate(ctx *config.Context, step *Step, meta *StepMetadata) (OutputFields, error) {

	ssmImport := *ctx.Project.Settings.SSMImport

	if !ssmImport.IsEnabled() {
		// step.Logger.Warn("SSM Prefix is blank. Will not attempt to load secrets.")
		step.Logger.Debugf("SSM Import is disabled. No secrets will be loaded.")
		return nil, nil
	}

	ssmPrefix, err := tmpl.New(ctx).Apply(ssmImport.GetPath())
	if err != nil {
		return nil, err
	}

	// Trim any trailing slash, then add our own
	ssmPrefix = strings.TrimSuffix(ssmPrefix, "/") + "/"

	logger := step.Logger.WithField("prefix", ssmPrefix)

	logger.Debugf("loading secrets from SSM")

	ssmClient := awsclients.SSMClient()

	request := &ssm.GetParametersByPathInput{
		Path:           aws.String(ssmPrefix),
		WithDecryption: aws.Bool(false), // dont need to decrypt, we just want the ARN
		Recursive:      ssmImport.Recursive,
	}

	paginator := ssm.NewGetParametersByPathPaginator(ssmClient, request, func(o *ssm.GetParametersByPathPaginatorOptions) {
		o.Limit = 10
	})

	secrets := make(map[string]config.EnvVar)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx.Context)
		if err != nil {
			logger.Warn("failed to load SSM Secrets")
			return nil, err
		}
		for _, parameter := range output.Parameters {
			// just want the last part of the name
			name := filepath.Base(*parameter.Name)
			secrets[name] = config.NewEnvVar(config.EnvVarTypeSSM, *parameter.ARN)
		}
	}

	ctx.Cache.SSMSecrets = secrets

	logger.Infof("loaded %d secrets from SSM", len(secrets))

	return nil, nil
}
