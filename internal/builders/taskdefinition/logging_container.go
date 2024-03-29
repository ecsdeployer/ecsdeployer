package taskdefinition

import (
	"errors"
	"strings"

	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

var ErrFirelensSSMUsageError = errors.New("Cannot use SSM references in firelens options")

// just the logging container, if desired. (firelens)
func (b *Builder) applyLoggingFirelensContainer() error {

	loggingConf := b.project.Logging

	if loggingConf.FirelensConfig.IsDisabled() {
		return nil
	}

	firelensConfig := loggingConf.FirelensConfig

	// flImageUri, err := b.tplEval(firelensConfig.Image.Value())
	// if err != nil {
	// 	return err
	// }
	// This is resolved already for us in the PreflightStep
	flImageUri, err := helpers.ResolveImageUri(b.ctx, firelensConfig.Image)
	if err != nil {
		return err
	}

	flConfig := &ecsTypes.FirelensConfiguration{
		Type: ecsTypes.FirelensConfigurationType(*firelensConfig.Type),
		// Options: make(map[string]string),
	}
	if firelensConfig.RouterOptions.HasSSM() {
		return ErrFirelensSSMUsageError
	}

	filteredRouterOpts := firelensConfig.RouterOptions.Filter()
	if len(filteredRouterOpts) > 0 {
		tpl := b.tpl()

		flConfig.Options = make(map[string]string)

		for lk, lv := range filteredRouterOpts {
			val, err := lv.GetValue(tpl)
			if err != nil {
				return err
			}
			flConfig.Options[lk] = val
		}
	}

	b.loggingContainer = &ecsTypes.ContainerDefinition{
		Name:                  firelensConfig.Name,
		Essential:             aws.Bool(true),
		Image:                 aws.String(flImageUri),
		FirelensConfiguration: flConfig,
	}

	if firelensConfig.Memory != nil {
		memVal := firelensConfig.Memory.GetValueOnly()
		if memVal > 0 {
			b.loggingContainer.MemoryReservation = aws.Int32(memVal)
		}
	}

	if firelensConfig.Credentials != nil {
		b.loggingContainer.RepositoryCredentials = &ecsTypes.RepositoryCredentials{
			CredentialsParameter: firelensConfig.Credentials,
		}
	}

	containerEnv := firelensConfig.EnvVars
	if *firelensConfig.InheritEnv {
		containerEnv = config.MergeEnvVarMaps(b.baseEnvVars, containerEnv)
	}
	if err := b.addEnvVarsToContainer(b.loggingContainer, containerEnv); err != nil {
		return err
	}

	if firelensConfig.LogToAwsLogs != nil && firelensConfig.LogToAwsLogs.Enabled() {

		awsLogOpts := config.EnvVarMap{
			// "awslogs-group":         config.NewEnvVar(config.EnvVarTypeTemplated, firelensConfig.LogToAwsLogs.Path),
			"awslogs-create-group":  config.NewEnvVar(config.EnvVarTypePlain, "true"),
			"awslogs-region":        config.NewEnvVar(config.EnvVarTypeTemplated, "{{ AwsRegion }}"),
			"awslogs-stream-prefix": config.NewEnvVar(config.EnvVarTypeTemplated, "firelens-{{.Name}}"),
		}

		logPathRaw, err := b.tplEval(firelensConfig.LogToAwsLogs.Path)
		if err != nil {
			return err
		}
		logPath, logStream, hasStream := strings.Cut(logPathRaw, ":")
		awsLogOpts["awslogs-group"] = config.NewEnvVar(config.EnvVarTypePlain, logPath)
		if hasStream {
			awsLogOpts["awslogs-stream-prefix"] = config.NewEnvVar(config.EnvVarTypePlain, logStream)
		}

		awsLogConfig := &config.TaskLoggingConfig{
			Driver:  util.Ptr(string(ecsTypes.LogDriverAwslogs)),
			Options: awsLogOpts,
		}

		if err := b.buildContainerLogging(b.loggingContainer, awsLogConfig); err != nil {
			return err
		}
	}

	return nil
}
