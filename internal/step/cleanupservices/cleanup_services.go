// this will delete services that are no longer used
package cleanupservices

import (
	"errors"
	"sync"
	"time"

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

type Step struct{}

func (Step) String() string {
	return "cleaning orphaned services"
}

func (Step) Skip(ctx *config.Context) bool {
	return !ctx.Project.Settings.KeepInSync.GetServices()
}

func (Step) Clean(ctx *config.Context) error {
	markerKey, markerVal, err := helpers.GetMarkerTag(ctx)
	if err != nil {
		log.WithError(err).Warn("Failed to determine marker tag. Unable to delete unused taskdefs")
		return step.Skip("Failed to determine marker tag. Unable to delete unused taskdefs")
	}

	expectedServiceNames := make([]string, 0, len(ctx.Project.Services))
	for _, svc := range ctx.Project.Services {
		svcName, err := tmpl.New(ctx).WithExtraFields(svc.TemplateFields()).Apply(*ctx.Project.Templates.ServiceName)
		if err != nil {
			log.WithError(err).Error("Unable to determine existing service names. Services will not be synced")
			return nil
		}
		expectedServiceNames = append(expectedServiceNames, svcName)
	}

	client := awsclients.TaggingClient()

	request := &tagging.GetResourcesInput{
		ResourceTypeFilters: []string{"ecs:service"},
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
			log.WithError(err).Warn("Failed to determine marker tag. Unable to delete unused services")
			return nil
		}
		for _, result := range output.ResourceTagMappingList {

			serviceName := helpers.GetECSServiceNameFromArn(*result.ResourceARN)
			if serviceName == "" {
				log.WithError(err).WithField("arn", result.ResourceARN).Warn("unparsable service ARN")
				continue
			}

			if slices.Contains(expectedServiceNames, serviceName) {
				// service is supposed to be there, so dont delete
				continue
			}

			wg.Add(1)

			go func(svcArn string) {
				defer wg.Done()
				err := destroyService(ctx, svcArn)
				if err != nil {
					log.WithError(err).WithField("service", svcArn).Warn("unable to delete. (ignoring)")
				}
			}(*result.ResourceARN)
		}
	}

	wg.Wait()

	return nil
}

func destroyService(ctx *config.Context, serviceArn string) error {

	serviceName := helpers.GetECSServiceNameFromArn(serviceArn)

	logger := log.WithField("service", serviceName)
	logger.Info("found unwanted service")

	clusterName := helpers.GetECSClusterNameFromArn(serviceArn)

	if clusterName == "" {
		clusterName = ctx.ClusterName()
	}

	client := awsclients.ECSClient()

	result, err := client.DescribeServices(ctx.Context, &ecs.DescribeServicesInput{
		Services: []string{serviceArn},
		Cluster:  &clusterName,
	})
	if err != nil {
		return err
	}

	if len(result.Services) == 0 {
		return errors.New("could not find service?")
	}

	svc := result.Services[0]

	// TODO: is this part even needed if we have force=true in delete?
	if svc.DesiredCount > 0 {
		// logger.Trace("Service has desired count greater than 1, spinning down")
		logger.Info("spinning service to zero...")
		_, err = client.UpdateService(ctx.Context, &ecs.UpdateServiceInput{
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
		Service: &serviceName,
		Cluster: &clusterName,
		Force:   aws.Bool(true),
	})

	if err != nil {
		return err
	}

	log.WithField("service", serviceName).Info("service deleted")

	return nil
}
