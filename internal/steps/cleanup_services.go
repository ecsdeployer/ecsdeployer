package steps

import (
	"errors"
	"sync"
	"time"

	"ecsdeployer.com/ecsdeployer/internal/awsclients"
	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	tagging "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	taggingTypes "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"github.com/caarlos0/log"
	"golang.org/x/exp/slices"
)

func CleanupServicesStep(resource *config.KeepInSync) *Step {

	if !*resource.Services {
		return NoopStep()
	}

	return NewStep(&Step{
		Label:    "CleanupServices",
		Resource: resource,
		Create:   stepCleanupServicesCreate,
	})
}

func stepCleanupServicesCreate(ctx *config.Context, step *Step, meta *StepMetadata) (OutputFields, error) {

	log := step.Logger

	markerTag, err := stepCleanupMarkerTag(ctx)
	if err != nil {
		log.WithError(err).Warn("Failed to determine marker tag. Unable to delete unused services")
		return nil, nil
	}

	expectedServiceNames := make([]string, 0, len(ctx.Project.Services))
	for _, svc := range ctx.Project.Services {
		svcName, err := tmpl.New(ctx).WithExtraFields(svc.TemplateFields()).Apply(*ctx.Project.Templates.ServiceName)
		if err != nil {
			log.WithError(err).Error("Unable to determine existing service names. Services will not be synced")
			return nil, nil
		}
		expectedServiceNames = append(expectedServiceNames, svcName)
	}

	client := awsclients.TaggingClient()

	request := &tagging.GetResourcesInput{
		ResourceTypeFilters: []string{"ecs:service"},
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
			log.WithError(err).Warn("Failed to determine marker tag. Unable to delete unused services")
			return nil, nil
		}
		for _, result := range output.ResourceTagMappingList {

			serviceName := helpers.GetECSServiceNameFromArn(*result.ResourceARN)
			if serviceName == "" {
				log.WithError(err).WithField("arn", result.ResourceARN).Warn("Received unparsable service ARN")
				continue
			}

			if slices.Contains(expectedServiceNames, serviceName) {
				// service is supposed to be there, so dont delete
				continue
			}

			// log.WithField("name", serviceName).Info("Found unwanted service. Will delete")

			wg.Add(1)

			go func(svcArn string) {
				defer wg.Done()
				err := destroyService(ctx, log, svcArn)
				if err != nil {
					log.WithError(err).WithField("service", svcArn).Warn("Unable to delete. Ignoring")
				}
			}(*result.ResourceARN)
		}
	}

	wg.Wait()

	return nil, nil
}

func destroyService(ctx *config.Context, log *log.Entry, serviceArn string) error {

	logger := log.WithField("service", serviceArn)

	clusterName := helpers.GetECSClusterNameFromArn(serviceArn)

	if clusterName == "" {
		cName, err := ctx.Project.Cluster.Name(ctx)
		if err != nil {
			return err
		}

		clusterName = cName
	}

	logger = logger.WithField("cluster", clusterName)

	logger.Info("Removing unused service")

	client := awsclients.ECSClient()

	result, err := client.DescribeServices(ctx.Context, &ecs.DescribeServicesInput{
		Services: []string{serviceArn},
		Cluster:  &clusterName,
		// Include:  []types.ServiceField{},
	})
	if err != nil {
		return err
	}

	if len(result.Services) == 0 {
		return errors.New("could not find service?")
	}

	svc := result.Services[0]

	if svc.DesiredCount > 0 {
		logger.Info("Service has desired count greater than 1, spinning down")
		_, err := client.UpdateService(ctx.Context, &ecs.UpdateServiceInput{
			Service:            &serviceArn,
			Cluster:            &clusterName,
			DesiredCount:       aws.Int32(0),
			ForceNewDeployment: true,
		})
		if err != nil {
			return err
		}

		// wait a lil bit for ecs to catch up
		if !helpers.IsTestingMode {
			time.Sleep(5 * time.Second)
		}
	}

	_, err = client.DeleteService(ctx.Context, &ecs.DeleteServiceInput{
		Service: &serviceArn,
		Cluster: &clusterName,
		Force:   aws.Bool(true),
	})

	if err != nil {
		return err
	}

	return nil
}
