package cronjob

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	cronBuilder "ecsdeployer.com/ecsdeployer/internal/builders/cron"
	"ecsdeployer.com/ecsdeployer/internal/substep/taskdefinition"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/webdestroya/go-log"
)

func (s *Step) runLegacyVariant(ctx *config.Context) error {
	return (&legacy{
		cronjob: s.cronjob,
	}).Run(ctx)
}

type legacy struct {
	cronjob    *config.CronJob
	taskDefArn string
	ruleArn    string
}

func (s *legacy) Run(ctx *config.Context) error {

	taskDefArn, err := taskdefinition.Register(ctx, s.cronjob)
	if err != nil {
		return err
	}
	s.taskDefArn = taskDefArn

	if err = s.putRule(ctx); err != nil {
		return err
	}

	if err = s.putTarget(ctx); err != nil {
		return err
	}

	return nil
}

func (s *legacy) putRule(ctx *config.Context) error {

	payload, err := cronBuilder.BuildCronRule(ctx, s.cronjob)
	if err != nil {
		return fmt.Errorf("failed to build PutRule: %w", err)
	}

	result, err := awsclients.EventsClient().PutRule(ctx, payload)
	if err != nil {
		return fmt.Errorf("failed to create Rule: %w", err)
	}

	s.ruleArn = *result.RuleArn

	return nil
}

func (s *legacy) putTarget(ctx *config.Context) error {
	putTargetsInput, err := cronBuilder.BuildCronTarget(ctx, s.cronjob, s.taskDefArn)
	if err != nil {
		return fmt.Errorf("failed to build PutTarget: %w", err)
	}

	client := awsclients.EventsClient()

	result, err := client.PutTargets(ctx.Context, putTargetsInput)
	if err != nil {
		return fmt.Errorf("failed to create Target: %w", err)
	}

	if result.FailedEntryCount > 0 {
		for _, failEntry := range result.FailedEntries {
			log.WithFields(log.Fields{
				"name":     s.cronjob.Name,
				"targetid": *failEntry.TargetId,
				"code":     *failEntry.ErrorCode,
				"message":  *failEntry.ErrorMessage,
			}).Errorf("failure")
		}
		return fmt.Errorf("Failed to create Cron Targets")
	}

	return nil
}
