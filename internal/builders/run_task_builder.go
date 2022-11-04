package builders

import (
	"errors"

	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func BuildRunTask(ctx *config.Context, resource *config.PreDeployTask) (*ecs.RunTaskInput, error) {

	project := ctx.Project
	taskDefaults := project.TaskDefaults
	templates := project.Templates

	tpl := tmpl.New(ctx).WithExtraFields(resource.TemplateFields())

	clusterArn, err := project.Cluster.Arn(ctx)
	if err != nil {
		return nil, err
	}

	platformVersion := util.Coalesce(resource.PlatformVersion, taskDefaults.PlatformVersion)

	runTaskInput := &ecs.RunTaskInput{
		TaskDefinition:           new(string),
		Cluster:                  aws.String(clusterArn),
		Count:                    aws.Int32(1),
		EnableECSManagedTags:     true,
		EnableExecuteCommand:     false,
		PlatformVersion:          platformVersion,
		PropagateTags:            ecsTypes.PropagateTagsTaskDefinition,
		CapacityProviderStrategy: config.NewSpotOnDemand().ExportCapacityStrategy(), // always ondemand
		// Overrides:            &ecstypes.TaskOverride{},
		// StartedBy:            new(string),
		// ReferenceId:          new(string),
		// Tags:                 []ecstypes.Tag{},
		// NetworkConfiguration: &ecstypes.NetworkConfiguration{},
		// PlacementConstraints:     []ecstypes.PlacementConstraint{},
		// PlacementStrategy:        []ecstypes.PlacementStrategy{},
		// LaunchType:               "",
	}

	/*
		if resource.Command != nil && len(*resource.Command) > 0 {

			// COMMAND OVERRIDE
			containerName, err := tpl.Apply(templates.ContainerName)
			if err != nil {
				return nil, err
			}
			containerOverride := ecstypes.ContainerOverride{
				Command: *resource.Command,
				Name:    aws.String(containerName),
			}

			runTaskInput.Overrides = &ecstypes.TaskOverride{
				ContainerOverrides: []ecstypes.ContainerOverride{containerOverride},
			}
		}
	*/

	// STARTED BY
	startedBy, err := tpl.Apply(*templates.PreDeployStartedBy)
	if err != nil {
		return nil, err
	}
	if startedBy != "" {
		runTaskInput.StartedBy = aws.String(startedBy)
	}

	// GROUP NAME
	groupName, err := tpl.Apply(*templates.PreDeployGroup)
	if err != nil {
		return nil, err
	}
	if groupName != "" {
		runTaskInput.Group = aws.String(groupName)
	}

	// NETWORK
	network := util.Coalesce(resource.Network, taskDefaults.Network, project.Network)
	if network == nil {
		return nil, errors.New("Unable to resolve network configuration!")
	}
	ecsNetworkConfig, err := network.ResolveECS(ctx)
	if err != nil {
		return nil, err
	}
	runTaskInput.NetworkConfiguration = ecsNetworkConfig

	return runTaskInput, nil
}
