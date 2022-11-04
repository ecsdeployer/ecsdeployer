package testutil

import (
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func LoadProjectConfig(file string) (*config.Context, error) {
	project, err := yaml.ParseYAMLFile[config.Project](file)
	if err != nil {
		return nil, err
	}

	return config.New(project), nil
}
