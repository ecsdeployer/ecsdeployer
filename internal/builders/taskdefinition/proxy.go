package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

func (b *Builder) applyProxyConfiguration() error {

	proxyConf := util.Coalesce(b.commonTask.ProxyConfig, b.taskDefaults.ProxyConfig)
	if proxyConf == nil {
		return nil
	}

	if proxyConf.Disabled {
		return nil
	}

	b.taskDef.ProxyConfiguration = &ecsTypes.ProxyConfiguration{}

	b.taskDef.ProxyConfiguration.Type = ecsTypes.ProxyConfigurationType(*proxyConf.Type)
	b.taskDef.ProxyConfiguration.ContainerName = proxyConf.ContainerName

	propList, _, err := config.ExportEnvVarMap(proxyConf.Properties, b.tpl(), envExportKVP, envExportKVP)
	if err != nil {
		return err
	}
	b.taskDef.ProxyConfiguration.Properties = propList

	return nil
}

func envExportKVP(k, v string) ecsTypes.KeyValuePair {
	return ecsTypes.KeyValuePair{
		Name:  &k,
		Value: &v,
	}
}
