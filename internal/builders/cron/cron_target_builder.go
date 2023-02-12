package cron

import (
	"encoding/json"
	"errors"

	"ecsdeployer.com/ecsdeployer/internal/tmpl"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	events "github.com/aws/aws-sdk-go-v2/service/eventbridge"
	eventTypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
)

func BuildCronTarget(ctx *config.Context, resource *config.CronJob, taskDefArn string) (*events.PutTargetsInput, error) {

	project := ctx.Project

	templates := project.Templates

	tpl := tmpl.New(ctx).WithExtraFields(resource.TemplateFields())

	// cronContainerName, err := tpl.Apply(templates.ContainerName)
	// if err != nil {
	// 	return nil, err
	// }

	ecsParams := &eventTypes.EcsParameters{
		TaskDefinitionArn:        aws.String(taskDefArn),
		TaskCount:                aws.Int32(1),
		EnableECSManagedTags:     true,
		EnableExecuteCommand:     false,
		LaunchType:               eventTypes.LaunchTypeFargate,
		PlatformVersion:          resource.PlatformVersion,
		PropagateTags:            eventTypes.PropagateTagsTaskDefinition,
		CapacityProviderStrategy: config.NewSpotOnDemand().ExportCapacityStrategyEventBridge(),
	}

	cronGroupName, err := tpl.Apply(*templates.CronGroup)
	if err != nil {
		return nil, err
	}
	if cronGroupName != "" {
		ecsParams.Group = aws.String(cronGroupName)
	}

	// Cronjob Input field
	// cronInput := cronInput{
	// 	ContainerOverrides: []cronContainerOverride{
	// 		{
	// 			Name:    cronContainerName,
	// 			Command: *resource.Command,
	// 		},
	// 	},
	// }

	// Because we have a unique taskdef for each cronjob, we dont need to override the input
	// this just makes it so an empty JSON is sent
	cronInput := cronInputObj{}
	cronInputJsonBytes, err := json.Marshal(cronInput)
	if err != nil {
		return nil, err
	}

	cronTargetName, err := tpl.Apply(*templates.CronTarget)
	if err != nil {
		return nil, err
	}

	cronRuleName, err := tpl.Apply(*templates.CronRule)
	if err != nil {
		return nil, err
	}

	clusterArn, err := project.Cluster.Arn(ctx)
	if err != nil {
		return nil, err
	}

	// Load network configuration
	network := util.Coalesce(resource.Network, project.TaskDefaults.Network, project.Network)
	if network == nil {
		return nil, errors.New("Unable to resolve network configuration!")
	}

	ecsNetworkConfig, err := network.ResolveCWE(ctx)
	if err != nil {
		return nil, err
	}
	ecsParams.NetworkConfiguration = ecsNetworkConfig

	// The target
	targetDef := eventTypes.Target{
		Arn:           aws.String(clusterArn),
		Id:            aws.String(cronTargetName),
		EcsParameters: ecsParams,
		Input:         aws.String(string(cronInputJsonBytes)),
		// DeadLetterConfig:            &eventtypes.DeadLetterConfig{},
	}
	if project.CronLauncherRole != nil {
		launcherRole, err := project.CronLauncherRole.Arn(ctx)
		if err != nil {
			return nil, err
		}
		targetDef.RoleArn = &launcherRole
	}

	// payload for aws call
	targetsDef := &events.PutTargetsInput{
		Rule:    aws.String(cronRuleName),
		Targets: []eventTypes.Target{targetDef},
	}

	if resource.EventBusName != nil {
		targetsDef.EventBusName = resource.EventBusName
	}

	return targetsDef, nil
}
