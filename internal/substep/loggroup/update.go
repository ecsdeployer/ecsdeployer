package loggroup

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/step"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	logs "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	logTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	log "github.com/caarlos0/log"
)

func (s *Substep) updateLogGroup(ctx *config.Context, current *logTypes.LogGroup) error {

	if !ctx.Project.Settings.KeepInSync.GetLogRetention() {
		return step.Skip("log retention sync disabled")
	}

	retention := ctx.Project.Logging.AwsLogConfig.Retention

	if current == nil {
		logGroup, err := s.describeLogGroup(ctx, true)
		if err != nil {
			return err
		}

		current = logGroup
	}

	if retention.EqualsLogGroup(*current) {
		// no updates needed
		return nil
	}

	if retention.Forever() {
		log.WithField("group", s.groupName).Debug("deleting log retention")
		_, err := awsclients.LogsClient().DeleteRetentionPolicy(ctx.Context, &logs.DeleteRetentionPolicyInput{
			LogGroupName: &s.groupName,
		})
		if err != nil {
			return err
		}
	}

	log.WithField("group", s.groupName).WithField("days", retention.Days()).Debug("updating log retention")
	return putRetentionPolicy(ctx, s.groupName, retention.Days())
}

func putRetentionPolicy(ctx *config.Context, logGroupName string, days int32) error {
	_, err := awsclients.LogsClient().PutRetentionPolicy(ctx.Context, &logs.PutRetentionPolicyInput{
		LogGroupName:    &logGroupName,
		RetentionInDays: &days,
	})
	if err != nil {
		return fmt.Errorf("failed to set retention for log %s: %w", logGroupName, err)
	}
	return nil
}
