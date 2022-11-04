package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func TestCommonTaskAttrs_Smoke(t *testing.T) {
	common := &config.CommonTaskAttrs{
		Architecture: util.Ptr(config.ArchitectureARM64),
		CommonContainerAttrs: config.CommonContainerAttrs{
			Name: "test",
		},
		Network: &config.NetworkConfiguration{
			Subnets: []config.NetworkFilter{
				{
					ID: util.Ptr("subnet-1111111"),
				},
			},
		},
	}

	if err := common.Validate(); err != nil {
		t.Errorf("Expected CommonTaskAttrs to be valid, but got: %s", err)
	}

	if !common.IsTaskStruct() {
		t.Errorf("Expected CommonTaskAttrs to pass IsTaskStruct()")
	}

	fields := common.TemplateFields()
	if fields["Name"] != "test" {
		t.Errorf("Expected TemplateFields to give Name of %s but got %s", "test", fields["Name"])
	}
	if fields["Arch"] != "arm64" {
		t.Errorf("Expected TemplateFields to give Arch of %s but got %s", "arm64", fields["Arch"])
	}
}
