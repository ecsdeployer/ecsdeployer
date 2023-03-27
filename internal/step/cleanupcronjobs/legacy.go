package cleanupcronjobs

import (
	"sync"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/deprecate"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/step"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	events "github.com/aws/aws-sdk-go-v2/service/eventbridge"
	tagging "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	taggingTypes "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	log "github.com/caarlos0/log"
	"golang.org/x/exp/slices"
)

func runLegacyCleanup(ctx *config.Context) error {
	deprecate.Deprecate_LegacyCron(ctx)

	markerKey, markerVal, err := helpers.GetMarkerTag(ctx)
	if err != nil {
		return step.Skipf("Failed to determine marker tag: %s", err.Error())
	}

	expectedRuleNames := make([]string, 0, len(ctx.Project.CronJobs))
	for _, cron := range ctx.Project.CronJobs {
		ruleName, err := tmpl.New(ctx).WithExtraFields(cron.TemplateFields()).Apply(*ctx.Project.Templates.CronRule)
		if err != nil {
			log.WithError(err).Error("Unable to determine existing cron names. Cronjobs will not be synced")
			return nil
		}
		expectedRuleNames = append(expectedRuleNames, ruleName)
	}

	client := awsclients.TaggingClient()

	request := &tagging.GetResourcesInput{
		ResourceTypeFilters: []string{"events:rule"},
		ResourcesPerPage:    aws.Int32(50),
		TagFilters: []taggingTypes.TagFilter{
			{
				Key:    aws.String(markerKey),
				Values: []string{markerVal},
			},
		},
	}

	paginator := tagging.NewGetResourcesPaginator(client, request)

	var wg sync.WaitGroup

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx.Context)
		if err != nil {
			log.WithError(err).Warn("Failed to page event:rule resources")
			return nil
		}
		for _, result := range output.ResourceTagMappingList {

			ruleName, _ := helpers.GetEventRuleNameAndBusFromArn(*result.ResourceARN)

			if slices.Contains(expectedRuleNames, ruleName) {
				// service is supposed to be there, so dont delete
				continue
			}

			wg.Add(1)

			go func(cronRuleArn string) {
				defer wg.Done()
				err := destroyRule(ctx, cronRuleArn)
				if err != nil {
					log.WithError(err).WithField("rule", cronRuleArn).Warn("Unable to delete. Ignoring")
				}
			}(*result.ResourceARN)
		}
	}

	wg.Wait()

	return nil
}

func destroyRule(ctx *config.Context, ruleArn string) error {

	logger := log.WithField("rule", ruleArn)

	ruleName, busName := helpers.GetEventRuleNameAndBusFromArn(ruleArn)

	logger.Info("Removing unused rule")

	client := awsclients.EventsClient()

	listTargetsReq := &events.ListTargetsByRuleInput{
		Rule: &ruleName,
	}
	if busName != "" {
		listTargetsReq.EventBusName = &busName
	}

	results, err := client.ListTargetsByRule(ctx.Context, listTargetsReq)
	if err != nil {
		return err
	}
	targetIds := make([]string, 0, len(results.Targets))
	for _, target := range results.Targets {
		targetIds = append(targetIds, *target.Id)
	}

	logger.Info("Removing targets from rule")
	removeTargetReq := &events.RemoveTargetsInput{
		Ids:  targetIds,
		Rule: &ruleName,
	}
	if busName != "" {
		removeTargetReq.EventBusName = &busName
	}
	_, err = client.RemoveTargets(ctx.Context, removeTargetReq)
	if err != nil {
		return err
	}

	deleteRuleReq := &events.DeleteRuleInput{
		Name: &ruleName,
	}
	if busName != "" {
		deleteRuleReq.EventBusName = &busName
	}

	_, err = client.DeleteRule(ctx.Context, deleteRuleReq)
	if err != nil {
		return err
	}

	return nil
}
