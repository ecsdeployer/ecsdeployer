package builders

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func BuildCreateService(ctx *config.Context, resource *config.Service) (*ecs.CreateServiceInput, error) {

	project := ctx.Project
	taskDefaults := project.TaskDefaults
	templates := project.Templates

	clusterArn, err := project.Cluster.Arn(ctx)
	if err != nil {
		return nil, err
	}

	commonTplFields := resource.TemplateFields()

	tpl := tmpl.New(ctx).WithExtraFields(commonTplFields)

	serviceName, err := tpl.Apply(*templates.ServiceName)
	if err != nil {
		return nil, err
	}

	// container name
	containerName, err := tpl.Apply(*templates.ContainerName)
	if err != nil {
		return nil, err
	}

	platformVersion := util.Coalesce(resource.PlatformVersion, taskDefaults.PlatformVersion, aws.String("LATEST"))

	createServiceInput := &ecs.CreateServiceInput{
		ServiceName:             &serviceName,
		Cluster:                 &clusterArn,
		DeploymentConfiguration: resource.RolloutConfig.GetAwsConfig(),
		DesiredCount:            aws.Int32(resource.DesiredCount),
		EnableECSManagedTags:    true,
		EnableExecuteCommand:    false,
		PlatformVersion:         platformVersion,
		PropagateTags:           ecsTypes.PropagateTagsTaskDefinition,
		TaskDefinition:          new(string),
		SchedulingStrategy:      ecsTypes.SchedulingStrategyReplica,
		// Tags:                    []ecstypes.Tag{},
		// NetworkConfiguration:    &ecstypes.NetworkConfiguration{},
		// LoadBalancers:           []ecstypes.LoadBalancer{},
		// HealthCheckGracePeriodSeconds: new(int32),
		// Role:                          new(string),
		// ClientToken:                   new(string),
		// LaunchType:                    "",
		// CapacityProviderStrategy:      []ecstypes.CapacityProviderStrategyItem{},
		// PlacementStrategy:             []ecstypes.PlacementStrategy{},
		// DeploymentController:          &ecstypes.DeploymentController{},
		// PlacementConstraints:          []ecstypes.PlacementConstraint{},
		// ServiceRegistries:             []ecstypes.ServiceRegistry{},
	}

	// eventually let the user specify the strategies
	spotOverride := util.Coalesce(resource.SpotOverride, taskDefaults.SpotOverride, &config.SpotOverrides{})
	createServiceInput.CapacityProviderStrategy = spotOverride.ExportCapacityStrategy()

	// NETWORK
	network := util.Coalesce(resource.Network, taskDefaults.Network, project.Network)
	if network == nil {
		return nil, errors.New("Unable to resolve network configuration!")
	}
	ecsNetworkConfig, err := network.ResolveECS(ctx)
	if err != nil {
		return nil, err
	}
	createServiceInput.NetworkConfiguration = ecsNetworkConfig

	if resource.IsLoadBalanced() {

		createServiceInput.LoadBalancers = make([]ecsTypes.LoadBalancer, 0, len(resource.LoadBalancers))

		for _, lbInfo := range resource.LoadBalancers {

			targetGroupArn, err := lbInfo.TargetGroup.Arn(ctx)
			if err != nil {
				return nil, err
			}

			createServiceInput.LoadBalancers = append(createServiceInput.LoadBalancers, ecsTypes.LoadBalancer{
				ContainerName:  &containerName,
				ContainerPort:  lbInfo.PortMapping.Port,
				TargetGroupArn: &targetGroupArn,
			})
		}

		createServiceInput.HealthCheckGracePeriodSeconds = resource.LoadBalancers.GetHealthCheckGracePeriod()
	}

	tagList, _, err := helpers.NameValuePair_Build_Tags(ctx, resource.Tags, commonTplFields, ecsTag)
	if err != nil {
		return nil, err
	}
	createServiceInput.Tags = tagList

	return createServiceInput, nil
}

func ecsTag(s1, s2 string) ecsTypes.Tag {
	return ecsTypes.Tag{
		Key:   &s1,
		Value: &s2,
	}
}
