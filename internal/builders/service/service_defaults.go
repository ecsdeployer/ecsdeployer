package service

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func (b *Builder) applyServiceDefaults() error {

	clusterArn, err := b.project.Cluster.Arn(b.ctx)
	if err != nil {
		return err
	}
	b.serviceDef.Cluster = &clusterArn

	serviceName, err := b.tplEval(*b.templates.ServiceName)
	if err != nil {
		return err
	}
	b.serviceDef.ServiceName = &serviceName

	// THIS IS PURPOSELY LEFT BLANK. IT WILL BE FILLED IN BY CALLER
	b.serviceDef.TaskDefinition = new(string)

	b.serviceDef.DeploymentConfiguration = b.entity.RolloutConfig.GetAwsConfig()

	// AZ rebalancing is incompatible with maximumPercent <= 100%.
	// When the rollout max is at or below 100%, explicitly disable it to avoid:
	// "InvalidParameterException: Availability Zone Rebalancing does not support maximumPercent <= 100%"
	if *b.entity.RolloutConfig.Maximum <= 100 {
		b.serviceDef.AvailabilityZoneRebalancing = ecsTypes.AvailabilityZoneRebalancingDisabled
	}

	b.serviceDef.DesiredCount = new(b.entity.DesiredCount)

	b.serviceDef.EnableECSManagedTags = true

	b.serviceDef.EnableExecuteCommand = false

	b.serviceDef.PropagateTags = ecsTypes.PropagateTagsTaskDefinition

	b.serviceDef.SchedulingStrategy = ecsTypes.SchedulingStrategyReplica

	platformVersion := util.Coalesce(b.entity.PlatformVersion, b.taskDefaults.PlatformVersion, new("LATEST"))
	b.serviceDef.PlatformVersion = platformVersion

	if b.project.ServiceRole != nil {
		serviceRoleArn, err := b.project.ServiceRole.Arn(b.ctx)
		if err != nil {
			return err
		}
		b.serviceDef.Role = &serviceRoleArn
	}

	return nil
}
