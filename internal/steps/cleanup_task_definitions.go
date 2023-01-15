package steps

import (
	"sync"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	tagging "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	taggingTypes "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"golang.org/x/exp/slices"
)

func CleanupTaskDefinitionsStep(resource *config.KeepInSync) *Step {

	if !*resource.TaskDefinitions {
		return NoopStep()
	}

	return NewStep(&Step{
		Label:    "CleanupTaskDefinitions",
		Resource: resource,
		Create:   stepCleanupTaskDefinitionsCreate,
	})
}

func stepCleanupTaskDefinitionsCreate(ctx *config.Context, step *Step, meta *StepMetadata) (OutputFields, error) {

	log := step.Logger

	markerTag, err := stepCleanupMarkerTag(ctx)
	if err != nil {
		log.WithError(err).Warn("Failed to determine marker tag. Unable to delete unused taskdefs")
		return nil, nil
	}

	expectedTaskFamilies := make([]string, 0, 20)

	consoleName, err := tmpl.New(ctx).WithExtraFields(ctx.Project.ConsoleTask.TemplateFields()).Apply(*ctx.Project.Templates.TaskFamily)
	if err != nil {
		log.WithError(err).Error("Unable to determine existing taskDef names. TaskDefs will not be purged")
		return nil, nil
	}
	expectedTaskFamilies = append(expectedTaskFamilies, consoleName)

	for _, svc := range ctx.Project.Services {
		svcName, err := tmpl.New(ctx).WithExtraFields(svc.TemplateFields()).Apply(*ctx.Project.Templates.TaskFamily)
		if err != nil {
			log.WithError(err).Error("Unable to determine existing taskDef names. TaskDefs will not be purged")
			return nil, nil
		}
		expectedTaskFamilies = append(expectedTaskFamilies, svcName)
	}

	for _, svc := range ctx.Project.PreDeployTasks {
		svcName, err := tmpl.New(ctx).WithExtraFields(svc.TemplateFields()).Apply(*ctx.Project.Templates.TaskFamily)
		if err != nil {
			log.WithError(err).Error("Unable to determine existing taskDef names. TaskDefs will not be purged")
			return nil, nil
		}
		expectedTaskFamilies = append(expectedTaskFamilies, svcName)
	}

	for _, svc := range ctx.Project.CronJobs {
		svcName, err := tmpl.New(ctx).WithExtraFields(svc.TemplateFields()).Apply(*ctx.Project.Templates.TaskFamily)
		if err != nil {
			log.WithError(err).Error("Unable to determine existing taskDef names. TaskDefs will not be purged")
			return nil, nil
		}
		expectedTaskFamilies = append(expectedTaskFamilies, svcName)
	}

	client := awsclients.TaggingClient()

	request := &tagging.GetResourcesInput{
		ResourceTypeFilters: []string{"ecs:task-definition"},
		ResourcesPerPage:    aws.Int32(50),
		TagFilters: []taggingTypes.TagFilter{
			{
				Key:    aws.String(markerTag.key),
				Values: []string{markerTag.value},
			},
		},
	}

	paginator := tagging.NewGetResourcesPaginator(client, request)

	var wg sync.WaitGroup

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx.Context)
		if err != nil {
			log.WithError(err).Warn("Failed to determine marker tag. Unable to delete unused taskDefs")
			return nil, nil
		}
		for _, result := range output.ResourceTagMappingList {

			taskFamilyName := helpers.GetTaskDefFamilyFromArn(*result.ResourceARN)
			if taskFamilyName == "" {
				log.WithError(err).WithField("arn", result.ResourceARN).Warn("Received unparsable taskDef ARN")
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
				log.WithField("taskDef", taskDefArn).Info("Deregistering")
				err := deregisterTaskDefinition(ctx, taskDefArn)
				if err != nil {
					log.WithError(err).WithField("taskDef", taskDefArn).Warn("Unable to delete. Ignoring")
				}
			}(*result.ResourceARN)
		}
	}

	wg.Wait()

	return nil, nil
}
