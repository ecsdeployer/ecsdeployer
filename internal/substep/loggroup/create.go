package loggroup

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	logs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	logTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/webdestroya/go-log"
)

func (s *Substep) createLogGroup(ctx *config.Context) error {

	common := s.entity.GetCommonTaskAttrs()

	tplFields := common.TemplateFields()

	_, tagMap, err := helpers.NameValuePair_Build_Tags[any](ctx, common.Tags, tplFields, nil)
	if err != nil {
		return err
	}

	logsClient := awsclients.LogsClient()

	request := &logs.CreateLogGroupInput{
		LogGroupName: &s.groupName,
		Tags:         tagMap,
	}

	log.WithField("group", s.groupName).Debug("creating log group")

	if _, err := logsClient.CreateLogGroup(ctx.Context, request); err != nil {
		var alreadyExistsErr *logTypes.ResourceAlreadyExistsException
		if !errors.As(err, &alreadyExistsErr) {
			return err
		}

		// it already exists? so pull it again and see
		return s.updateLogGroup(ctx, nil)
	}

	retention := ctx.Project.Logging.AwsLogConfig.Retention
	if !retention.Forever() {
		return putRetentionPolicy(ctx, s.groupName, retention.Days())
	}

	return nil
}
