package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// parent function to append any sidecars
func (b *Builder) applySidecarContainers() error {

	sidecars := util.CoalesceArray(b.commonTask.Sidecars, b.taskDefaults.Sidecars)

	if sidecars == nil || len(sidecars) == 0 {
		return nil
	}

	for _, sidecar := range sidecars {
		if err := b.applySidecarContainer(sidecar); err != nil {
			return err
		}
	}

	return nil
}

func (b *Builder) applySidecarContainer(sidecar *config.Sidecar) error {

	cdef := &ecsTypes.ContainerDefinition{}

	if err := b.applyContainerDefaults(cdef, sidecar); err != nil {
		return err
	}

	if sidecar.InheritEnv {
		if err := b.addEnvVarsToContainer(cdef, b.baseEnvVars); err != nil {
			return err
		}
	}

	if err := b.applyContainerLogging(cdef, sidecar); err != nil {
		return err
	}

	// AT THE VERY END
	b.sidecars = append(b.sidecars, cdef)
	return nil
}
