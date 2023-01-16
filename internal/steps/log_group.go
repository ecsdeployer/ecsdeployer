package steps

import (
	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	logs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	logTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

// type LogGroupStepParams struct {
// 	LogGroupName string
// }

func LogGroupStep(resource interface{}) *Step {
	common, err := config.ExtractCommonTaskAttrs(resource)
	if err != nil {
		panic(err)
	}

	return NewStep(&Step{
		Label:    "LogGroup",
		ID:       common.Name,
		Resource: common,
		Create:   stepLogGroupCreate,
		Read:     stepLogGroupRead,
		PreApply: stepLogGroupPreApply,
		Update:   stepLogGroupUpdate,
	})
}

func stepLogGroupPreApply(ctx *config.Context, step *Step, meta *StepMetadata) error {

	if ctx.Project.Logging.IsDisabled() || ctx.Project.Logging.AwsLogConfig.IsDisabled() {
		return nil
	}

	common := (step.Resource).(*config.CommonTaskAttrs)

	if common.LoggingConfig != nil {
		step.Logger.Debug("AwsLogs have been disabled for this task")
		// logging disabled
		return nil
	}

	tpl := tmpl.New(ctx).WithExtraFields(common.TemplateFields())

	logGroupname, err := tpl.Apply(*ctx.Project.Templates.LogGroup)
	if err != nil {
		return err
	}

	step.Attributes["logGroupName"] = logGroupname

	return nil
}

func stepLogGroupCreate(ctx *config.Context, step *Step, meta *StepMetadata) (OutputFields, error) {

	logGroupNameRes, ok := step.Attributes["logGroupName"]
	if !ok {
		// the attribute key is missing, that means this should be skipped
		return nil, nil // errors.New("LogGroupName key is missing, skip stepLogGroupCreate")
	}
	logGroupName := logGroupNameRes.(string)

	logger := step.Logger.WithField("logGroup", logGroupName)

	outputs := make(OutputFields, 1)
	outputs["LogGroupName"] = logGroupName

	logsClient := awsclients.LogsClient()

	request := &logs.CreateLogGroupInput{
		LogGroupName: aws.String(logGroupName),
		Tags:         make(map[string]string),
		// KmsKeyId:     new(string),
	}

	common := (step.Resource).(*config.CommonTaskAttrs)

	commonTpl, err := helpers.GetDefaultTaskTemplateFields(ctx, common)
	if err != nil {
		return nil, err
	}

	_, tagMap, err := helpers.NameValuePair_Build_Tags[interface{}](ctx, common.Tags, commonTpl, nil)
	if err != nil {
		return nil, err
	}
	request.Tags = tagMap

	// CREATE LOG GROUP
	_, err = logsClient.CreateLogGroup(ctx.Context, request)
	if err != nil {
		return nil, err
	}
	logger.Info("Created Log Group")

	if !ctx.Project.Logging.AwsLogConfig.Retention.Forever() {
		// PUT RETENTION SETTINGS
		_, err = logsClient.PutRetentionPolicy(ctx.Context, &logs.PutRetentionPolicyInput{
			LogGroupName:    aws.String(logGroupName),
			RetentionInDays: ctx.Project.Logging.AwsLogConfig.Retention.ToAwsInt32(),
		})
		if err != nil {
			// logger.WithError(err).Warn("Failed to set retention for log group")
			// return nil, errors.New("Failed to set retention for log group")
			return nil, err
		}
	}

	return outputs, nil
}

func stepLogGroupRead(ctx *config.Context, step *Step, meta *StepMetadata) (any, error) {

	logGroupName, ok := step.Attributes["logGroupName"]
	if !ok {
		// the attribute key is missing, that means this should be skipped
		return aws.Bool(true), nil
	}

	for _, lg := range ctx.Cache.LogGroups {
		if *lg.LogGroupName == logGroupName {
			return lg, nil
		}
	}

	return nil, nil
}

func stepLogGroupUpdate(ctx *config.Context, step *Step, meta *StepMetadata) (OutputFields, error) {

	logGroupNameRes, ok := step.Attributes["logGroupName"]
	if !ok {
		// the attribute key is missing, that means this should be skipped
		return nil, nil // errors.New("logGroupUpdate - key is missing")
	}
	logGroupName := logGroupNameRes.(string)

	logger := step.Logger.WithField("logGroup", logGroupName)

	if !*ctx.Project.Settings.KeepInSync.LogRetention {
		logger.Info("Log Retention Sync disabled. Skipping")
		return nil, nil
	}

	if step.ExistingResource == nil {
		logger.Warn("LogGroup update is missing the existing resource??")
		return nil, nil
	}

	// outputs := make(OutputFields, 1)
	// outputs["LogGroupName"] = logGroupName

	logConfig := ctx.Project.Logging.AwsLogConfig.Retention

	logsClient := awsclients.LogsClient()

	logGroup := step.ExistingResource.(logTypes.LogGroup)

	if logConfig.Forever() {
		if logGroup.RetentionInDays == nil {
			// all good
			return nil, nil
		}

		// want infinite retention, but it has a retention
		_, err := logsClient.DeleteRetentionPolicy(ctx.Context, &logs.DeleteRetentionPolicyInput{
			LogGroupName: aws.String(logGroupName),
		})
		if err != nil {
			logger.Warn("Failed to remove retention policy for log group")
			return nil, err
		}
		logger.Info("Removed retention policy for log group")
		return nil, nil
	}

	if logConfig.EqualsLogGroup(logGroup) {
		// all good, nothing to fix
		return nil, nil
	}

	_, err := logsClient.PutRetentionPolicy(ctx.Context, &logs.PutRetentionPolicyInput{
		LogGroupName:    aws.String(logGroupName),
		RetentionInDays: logConfig.ToAwsInt32(),
	})
	if err != nil {
		logger.Warn("Failed to set retention for log group")
		return nil, err
	}
	return nil, nil
}
