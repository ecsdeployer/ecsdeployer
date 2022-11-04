package fargate

import (
	"fmt"
)

type FargateResource struct {
	Cpu    int32 `json:"cpu"`
	Memory int32 `json:"memory"`
}

func (fr *FargateResource) Fits(cpu int32, memory int32) bool {
	return fr.Cpu >= cpu && fr.Memory >= memory
}

func (fr *FargateResource) CpuString() string {
	return fmt.Sprintf("%d", fr.Cpu)
}

func (fr *FargateResource) MemoryString() string {
	return fmt.Sprintf("%d", fr.Memory)
}
