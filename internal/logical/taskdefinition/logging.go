package taskdefinition

import (
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
)

type LogConfig struct {
	Driver ecsTypes.LogDriver

	Options config.EnvVarMap
}

func (obj *LogConfig) ExportForContainer(cont *Container) (*ecsTypes.LogConfiguration, error) {
	conf := &ecsTypes.LogConfiguration{
		LogDriver:     "",
		Options:       map[string]string{},
		SecretOptions: []ecsTypes.Secret{},
	}
	return conf, nil
}

type FirelensConfig struct {
	Type ecsTypes.FirelensConfigurationType

	Options config.EnvVarMap
}

func (obj *FirelensConfig) ExportForContainer(cont *Container) (*ecsTypes.FirelensConfiguration, error) {
	conf := &ecsTypes.FirelensConfiguration{
		Type:    "",
		Options: map[string]string{},
	}
	return conf, nil
}
