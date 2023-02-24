package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func (b *Builder) applyTaskDefaults() error {

	familyName, err := b.tplEval(*b.templates.TaskFamily)
	if err != nil {
		return err
	}
	b.taskDef.Family = aws.String(familyName)

	b.taskDef.NetworkMode = ecsTypes.NetworkModeAwsvpc

	b.taskDef.ContainerDefinitions = make([]ecsTypes.ContainerDefinition, 0)

	b.taskDef.RequiresCompatibilities = []ecsTypes.Compatibility{
		ecsTypes.CompatibilityFargate,
	}

	arch := util.Coalesce(b.commonTask.Architecture, b.taskDefaults.Architecture, util.Ptr(config.ArchitectureDefault))

	b.taskDef.RuntimePlatform = &ecsTypes.RuntimePlatform{
		CpuArchitecture:       arch.ToAws(),
		OperatingSystemFamily: ecsTypes.OSFamilyLinux,
	}

	return nil
}
