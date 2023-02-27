package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func ecsKeyValuePair(k, v string) ecsTypes.KeyValuePair {
	return ecsTypes.KeyValuePair{
		Name:  &k,
		Value: &v,
	}
}

func ecsSecret(k, v string) ecsTypes.Secret {
	return ecsTypes.Secret{
		Name:      &k,
		ValueFrom: &v,
	}
}

func (b *Builder) createDeploymentEnvVars() error {
	b.deploymentEnvVars = make(config.EnvVarMap)

	if b.project.Settings.SkipDeploymentEnvVars {
		return nil
	}
	// add the deployment env vars
	for k, v := range config.DefaultDeploymentEnvVars {
		b.deploymentEnvVars[k] = config.NewEnvVar(config.EnvVarTypeTemplated, v)
	}

	return nil
}

// creates the env var map used as the baseline
func (b *Builder) createTaskEnvVars() error {
	b.baseEnvVars = config.MergeEnvVarMaps(make(config.EnvVarMap), b.deploymentEnvVars)

	if len(b.ctx.Cache.SSMSecrets) > 0 {
		b.baseEnvVars = config.MergeEnvVarMaps(b.baseEnvVars, b.ctx.Cache.SSMSecrets)
	}
	b.baseEnvVars = config.MergeEnvVarMaps(b.baseEnvVars, b.taskDefaults.EnvVars)
	b.baseEnvVars = config.MergeEnvVarMaps(b.baseEnvVars, b.commonTask.EnvVars)
	return nil
}

// func (b *Builder) buildEnvVarMap(newVars config.EnvVarMap, inheritEnv bool) config.EnvVarMap {
// 	if newVars == nil {
// 		newVars = make(config.EnvVarMap)
// 	}

// 	if inheritEnv {
// 		newVars = config.MergeEnvVarMaps(b.baseEnvVars, newVars)
// 	}

// 	return newVars
// }

func (b *Builder) addEnvVarsToContainer(cdef *ecsTypes.ContainerDefinition, varMap config.EnvVarMap) error {

	envTpl := b.containerTpl(cdef)

	envvars, secrets, err := config.ExportEnvVarMap(varMap, envTpl, ecsKeyValuePair, ecsSecret)
	if err != nil {
		return err
	}

	cdef.Environment = envvars
	cdef.Secrets = secrets

	return nil
}
