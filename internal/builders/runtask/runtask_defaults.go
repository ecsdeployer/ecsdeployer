package runtask

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func (b *Builder) applyRunTaskDefaults() error {

	clusterArn, err := b.project.Cluster.Arn(b.ctx)
	if err != nil {
		return err
	}
	b.runTaskDef.Cluster = &clusterArn

	// THIS IS PURPOSELY LEFT BLANK. IT WILL BE FILLED IN BY CALLER
	b.runTaskDef.TaskDefinition = new(string)

	b.runTaskDef.Count = aws.Int32(1)

	b.runTaskDef.EnableECSManagedTags = true

	b.runTaskDef.EnableExecuteCommand = false

	b.runTaskDef.PropagateTags = ecsTypes.PropagateTagsTaskDefinition

	platformVersion := util.Coalesce(b.entity.PlatformVersion, b.taskDefaults.PlatformVersion, aws.String("LATEST"))
	b.runTaskDef.PlatformVersion = platformVersion

	// STARTED BY
	startedBy, err := b.tplEval(*b.templates.PreDeployStartedBy)
	if err != nil {
		return err
	}
	if !util.IsBlank(&startedBy) {
		b.runTaskDef.StartedBy = &startedBy
	}

	// GROUP NAME
	groupName, err := b.tplEval(*b.templates.PreDeployGroup)
	if err != nil {
		return err
	}
	if !util.IsBlank(&groupName) {
		b.runTaskDef.Group = &groupName
	}

	return nil
}
