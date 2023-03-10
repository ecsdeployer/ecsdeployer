package service

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

func BuildUpdate(ctx *config.Context, entity *config.Service) (*ecs.UpdateServiceInput, error) {
	createServiceInput, err := BuildCreate(ctx, entity)
	if err != nil {
		return nil, err
	}

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
