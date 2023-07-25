package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

// parent function to append any sidecars
func (b *Builder) applySidecarContainers() error {

	// do not merge from task defaults. if they specify sidecars, then that is the only sidecar list
	sidecars := util.CoalesceArray(b.commonTask.Sidecars, b.taskDefaults.Sidecars)

	if len(sidecars) == 0 {
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

	cdef.Essential = sidecar.Essential

	scEnvVars := make(config.EnvVarMap, 0)

	if sidecar.InheritEnv {
		scEnvVars = config.MergeEnvVarMaps(scEnvVars, b.baseEnvVars)
	}

	// if they wanted env vars, add them (overriding anything that was inherited)
	if len(sidecar.EnvVars) > 0 {
		scEnvVars = config.MergeEnvVarMaps(scEnvVars, sidecar.EnvVars)
	}

	if len(scEnvVars) > 0 {
		if err := b.addEnvVarsToContainer(cdef, scEnvVars); err != nil {
			return err
		}
	}

	if err := b.applySidecarPortMappings(cdef, sidecar); err != nil {
		return err
	}
	if err := b.applyContainerLogging(cdef, sidecar); err != nil {
		return err
	}

	// AT THE VERY END
	b.sidecars = append(b.sidecars, cdef)
	return nil
}
