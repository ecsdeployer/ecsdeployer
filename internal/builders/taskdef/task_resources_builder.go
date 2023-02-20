package taskdef

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/builders/pipeline"
	"ecsdeployer.com/ecsdeployer/internal/fargate"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type TaskResourcesBuilder struct {
	Resource config.IsTaskStruct
}

func (pc *TaskResourcesBuilder) Apply(obj *pipeline.PipeItem[ecs.RegisterTaskDefinitionInput]) error {

	common, err := config.ExtractCommonTaskAttrs(pc.Resource)
	if err != nil {
		return err
	}

	project := obj.Context.Project
	taskDefaults := project.TaskDefaults

	storage := util.Coalesce(common.Storage, taskDefaults.Storage)
	if storage != nil {
		obj.Data.EphemeralStorage = &ecsTypes.EphemeralStorage{
			SizeInGiB: int32(*storage),
		}
	}

	// select fargate resources
	cpu := util.Coalesce(common.Cpu, taskDefaults.Cpu)
	memory := util.Coalesce(common.Memory, taskDefaults.Memory)
	if cpu == nil || memory == nil {
		return fmt.Errorf("You need to specify the CPU/Memory on the task defaults")
	}
	memoryValue, err := memory.MegabytesFromCpu(cpu)
	if err != nil {
		return err
	}

	fargateResource := fargate.FindFargateBestFit(cpu.Shares(), memoryValue)
	obj.Data.Cpu = aws.String(fargateResource.CpuString())
	obj.Data.Memory = aws.String(fargateResource.MemoryString())

	return nil
}
