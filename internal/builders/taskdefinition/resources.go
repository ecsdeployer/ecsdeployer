package taskdefinition

import (
	"fmt"

	"ecsdeployer.com/ecsdeployer/internal/fargate"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func (b *Builder) applyTaskResources() error {
	storage := util.Coalesce(b.commonTask.Storage, b.taskDefaults.Storage)
	if storage != nil {
		b.taskDef.EphemeralStorage = &ecsTypes.EphemeralStorage{
			SizeInGiB: storage.Gigabytes(),
		}
	}

	// select fargate resources
	cpu := util.Coalesce(b.commonTask.Cpu, b.taskDefaults.Cpu)
	memory := util.Coalesce(b.commonTask.Memory, b.taskDefaults.Memory)
	if cpu == nil || memory == nil {
		return fmt.Errorf("You need to specify the CPU/Memory on the task defaults")
	}
	memoryValue, err := memory.MegabytesFromCpu(cpu)
	if err != nil {
		return err
	}

	fargateResource := fargate.FindFargateBestFit(cpu.Shares(), memoryValue)
	b.taskDef.Cpu = aws.String(fargateResource.CpuString())
	b.taskDef.Memory = aws.String(fargateResource.MemoryString())

	return nil
}

// this is meant for containers OTHER than the primary
func (b *Builder) applyContainerResources(cdef *ecsTypes.ContainerDefinition, thing hasContainerAttrs) error {

	common := thing.GetCommonContainerAttrs()

	if common.Cpu != nil {
		cdef.Cpu = common.Cpu.Shares()
	}

	if common.Memory != nil {
		memoryValue, err := common.Memory.MegabytesFromCpu(common.Cpu)
		if err != nil {
			return err
		}

		if memoryValue > 0 {
			cdef.MemoryReservation = aws.Int32(memoryValue)
		}
	}

	return nil
}
