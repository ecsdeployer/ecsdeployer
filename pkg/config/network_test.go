package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

type networkTestStruct struct {
	Network *config.NetworkConfiguration `yaml:"network"`
}

func TestNetwork_Valid(t *testing.T) {
	// network := config.NetworkConfiguration{}

	// obj, err := ParseYAMLFile("testdata/network/full.yml", networkTestStruct{})
	_, err := yaml.ParseYAMLFile[networkTestStruct]("testdata/network/combined.yml")
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}

}
