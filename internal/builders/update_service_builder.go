package builders

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func BuildUpdateService(ctx *config.Context, resource *config.Service) (*ecs.UpdateServiceInput, error) {

	createServiceInput, err := BuildCreateService(ctx, resource)
	if err != nil {
		return nil, err
	}

	// just pull all the info from the CreateService call... it's all the same
	updateServiceInput := &ecs.UpdateServiceInput{
		Service:                       createServiceInput.ServiceName,
		CapacityProviderStrategy:      createServiceInput.CapacityProviderStrategy,
		Cluster:                       createServiceInput.Cluster,
		DeploymentConfiguration:       createServiceInput.DeploymentConfiguration,
		DesiredCount:                  createServiceInput.DesiredCount,
		EnableECSManagedTags:          aws.Bool(createServiceInput.EnableECSManagedTags),
		EnableExecuteCommand:          aws.Bool(createServiceInput.EnableExecuteCommand),
		ForceNewDeployment:            true,
		HealthCheckGracePeriodSeconds: createServiceInput.HealthCheckGracePeriodSeconds,
		LoadBalancers:                 createServiceInput.LoadBalancers,
		NetworkConfiguration:          createServiceInput.NetworkConfiguration,
		PlacementConstraints:          createServiceInput.PlacementConstraints,
		PlacementStrategy:             createServiceInput.PlacementStrategy,
		PlatformVersion:               createServiceInput.PlatformVersion,
		PropagateTags:                 createServiceInput.PropagateTags,
		ServiceRegistries:             createServiceInput.ServiceRegistries,
		TaskDefinition:                createServiceInput.TaskDefinition,
	}

	return updateServiceInput, nil
}
