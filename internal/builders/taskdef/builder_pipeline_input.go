package taskdef

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
)

type pipelineInput struct {
	TaskDef  *ecs.RegisterTaskDefinitionInput
	Context  *config.Context
	Common   *config.CommonTaskAttrs
	Resource config.IsTaskStruct
}

type TaskDefPipelineApplierFunc func(input *pipelineInput) error
