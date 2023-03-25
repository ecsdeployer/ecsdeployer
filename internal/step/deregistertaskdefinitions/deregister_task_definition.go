// This deregisters task definitions that are still being managed by ecsdeployer
package deregistertaskdefinitions

import (
	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/semerrgroup"
	"ecsdeployer.com/ecsdeployer/internal/step"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	log "github.com/caarlos0/log"
	"golang.org/x/exp/slices"
)

type Step struct{}

func (Step) String() string {
	return "cleaning previous task definitions"
}

func (Step) Skip(ctx *config.Context) bool {
	return ctx.Project.Settings.KeepInSync.GetTaskDefinitions() != config.KeepInSyncTaskDefinitionsEnabled &&
		ctx.Project.Settings.KeepInSync.GetTaskDefinitions() != config.KeepInSyncTaskDefinitionsOnlyManaged
}

func (Step) Clean(ctx *config.Context) error {
	g := semerrgroup.NewSkipAware(semerrgroup.New(ctx.Concurrency(5)))
	for _, defArn := range ctx.Cache.TaskDefinitions() {
		family := helpers.GetTaskDefFamilyFromArn(defArn)
		g.Go(func() error {
			if err := deregisterPreviousTaskFamily(ctx, family); err != nil {
				if step.IsSkip(err) {
					return nil
				}
				log.WithField("reason", err.Error()).WithField("family", family).Warn("failed to deregister task definition")
			}
			return nil
		})
	}

	return g.Wait()
}

func deregisterPreviousTaskFamily(ctx *config.Context, family string) error {
	ecsClient := awsclients.ECSClient()

	request := &ecs.ListTaskDefinitionsInput{
		FamilyPrefix: &family,
		Sort:         ecsTypes.SortOrderDesc,
		Status:       ecsTypes.TaskDefinitionStatusActive,
	}

	paginator := ecs.NewListTaskDefinitionsPaginator(ecsClient, request)

	oldTaskDefs := make([]string, 0, 10)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx.Context)
		if err != nil {
			return step.Skipf("Failed to list task definitions: %s", err)
			// return err
		}

		oldTaskDefs = append(oldTaskDefs, output.TaskDefinitionArns...)
	}

	if len(oldTaskDefs) == 0 {
		// this should not happen... we just registered a task def earlier, how could we not have any
		log.WithField("family", family).Warn("no task defs found???")
		// logger.Debug("No task definitions found.")
		return nil
	}

	if len(oldTaskDefs) == 1 {
		log.WithField("family", family).Trace("nothing to deregister")
		// logger.Debug("Only 1 task definition found.")
		return nil
	}

	for _, oldArn := range oldTaskDefs {

		if slices.Contains(ctx.Cache.TaskDefinitions(), oldArn) {
			continue
		}

		_, err := ecsClient.DeregisterTaskDefinition(ctx.Context, &ecs.DeregisterTaskDefinitionInput{
			TaskDefinition: &oldArn,
		})
		if err != nil {
			log.WithField("arn", oldArn).WithError(err).Warn("deregistering failed")
		} else {
			log.WithField("arn", oldArn).Trace("deregistered")
		}

	}
	return nil
}
