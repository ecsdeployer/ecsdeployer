package preloadsecrets

import (
	"fmt"
	"path/filepath"
	"strings"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/webdestroya/go-log"
)

type Step struct{}

func (Step) String() string {
	return "preloading ssm secrets"
}

func (Step) Skip(ctx *config.Context) bool {
	return !ctx.Project.Settings.SSMImport.IsEnabled()
}

func (Step) Preload(ctx *config.Context) error {

	ssmImport := *ctx.Project.Settings.SSMImport

	ssmPrefix, err := tmpl.New(ctx).Apply(ssmImport.GetPath())
	if err != nil {
		return err
	}

	// Trim any trailing slash, then add our own
	ssmPrefix = strings.TrimSuffix(ssmPrefix, "/") + "/"

	// logger := step.Logger.WithField("prefix", ssmPrefix)

	// logger.Debug("loading secrets from SSM")

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
			// logger.Warn("failed to load SSM Secrets")
			return fmt.Errorf("failed to load secrets: %w", err)
		}
		for _, parameter := range output.Parameters {
			// just want the last part of the name
			name := filepath.Base(*parameter.Name)
			secrets[name] = config.NewEnvVar(config.EnvVarTypeSSM, *parameter.ARN)
		}
	}

	ctx.Cache.SSMSecrets = secrets
	ctx.Cache.SSMSecretsCached = true

	log.WithField("total", len(ctx.Cache.SSMSecrets)).WithField("prefix", ssmPrefix).Info("imported secrets")

	return nil
}
