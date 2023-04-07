// this will deregister task definitions that are no longer being managed by ECSDeployer
package cleanuptaskdefinitions

import (
	"sync"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/step"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	tagging "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	taggingTypes "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"github.com/webdestroya/go-log"
	"golang.org/x/exp/slices"
)

var ErrUnableToDetermineTaskDefsError = step.Skip("Unable to determine existing taskDef names. TaskDefs will not be purged")

type Step struct{}

func (Step) String() string {
	return "cleaning orphaned task definitions"
}

func (Step) Skip(ctx *config.Context) bool {
	return ctx.Project.Settings.KeepInSync.GetTaskDefinitions() != config.KeepInSyncTaskDefinitionsEnabled
}

func (Step) Clean(ctx *config.Context) error {

	markerKey, markerVal, err := helpers.GetMarkerTag(ctx)
	if err != nil {
		log.WithError(err).Warn("Failed to determine marker tag. Unable to delete unused taskdefs")
		return nil
	}

	expectedTaskFamilies := make([]string, 0, 20)

	consoleName, err := tmpl.New(ctx).WithExtraFields(ctx.Project.ConsoleTask.TemplateFields()).Apply(*ctx.Project.Templates.TaskFamily)
	if err != nil {
		return ErrUnableToDetermineTaskDefsError
	}
	expectedTaskFamilies = append(expectedTaskFamilies, consoleName)

	for _, svc := range ctx.Project.Services {
		svcName, err := tmpl.New(ctx).WithExtraFields(svc.TemplateFields()).Apply(*ctx.Project.Templates.TaskFamily)
		if err != nil {
			return ErrUnableToDetermineTaskDefsError
		}
		expectedTaskFamilies = append(expectedTaskFamilies, svcName)
	}

	for _, svc := range ctx.Project.PreDeployTasks {
		svcName, err := tmpl.New(ctx).WithExtraFields(svc.TemplateFields()).Apply(*ctx.Project.Templates.TaskFamily)
		if err != nil {
			return ErrUnableToDetermineTaskDefsError
		}
		expectedTaskFamilies = append(expectedTaskFamilies, svcName)
	}

	for _, svc := range ctx.Project.CronJobs {
		svcName, err := tmpl.New(ctx).WithExtraFields(svc.TemplateFields()).Apply(*ctx.Project.Templates.TaskFamily)
		if err != nil {
			return ErrUnableToDetermineTaskDefsError
		}
		expectedTaskFamilies = append(expectedTaskFamilies, svcName)
	}

	client := awsclients.TaggingClient()

	request := &tagging.GetResourcesInput{
		ResourceTypeFilters: []string{"ecs:task-definition"},
		ResourcesPerPage:    aws.Int32(50),
		TagFilters: []taggingTypes.TagFilter{
			{
				Key:    &markerKey,
				Values: []string{markerVal},
			},
		},
	}

	paginator := tagging.NewGetResourcesPaginator(client, request)

	var wg sync.WaitGroup

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx.Context)
		if err != nil {
			log.WithError(err).Warn("Failed to determine marker tag. Unable to delete unused taskDefs")
			return nil
		}
		for _, result := range output.ResourceTagMappingList {

			taskFamilyName := helpers.GetTaskDefFamilyFromArn(*result.ResourceARN)
			if taskFamilyName == "" {
				log.WithError(err).WithField("arn", *result.ResourceARN).Warn("Received unparsable taskDef ARN")
				continue
			}

			if slices.Contains(expectedTaskFamilies, taskFamilyName) {
				// service is supposed to be there, so dont delete
				continue
			}

			// log.WithField("name", serviceName).Info("Found unwanted service. Will delete")

			wg.Add(1)

			go func(taskDefArn string) {
				defer wg.Done()
				_ = DeregisterTaskDefinition(ctx, taskDefArn)
			}(*result.ResourceARN)
		}
	}

	wg.Wait()

	return nil
}

func DeregisterTaskDefinition(ctx *config.Context, taskDefinitionArn string) error {
	ecsClient := awsclients.ECSClient()

	_, err := ecsClient.DeregisterTaskDefinition(ctx.Context, &ecs.DeregisterTaskDefinitionInput{
		TaskDefinition: &taskDefinitionArn,
	})
	if err != nil {
		log.WithError(err).WithField("arn", taskDefinitionArn).Warn("deregistration failed. ignoring")
		return err
	}
	log.WithField("arn", taskDefinitionArn).Debug("deregistering task definition")

	return nil
}
