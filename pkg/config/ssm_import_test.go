package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func TestSSMImport_Unmarshal(t *testing.T) {

	defSSM := config.SSMImport{}
	defSSM.ApplyDefaults()

	bTrue := true
	// bFalse := false

	tables := []struct {
		str      string
		expected config.SSMImport
	}{
		{"ssm_import: true", config.SSMImport{Enabled: true, Path: defSSM.Path, Recursive: defSSM.Recursive}},
		{"ssm_import: false", config.SSMImport{Enabled: false, Path: defSSM.Path, Recursive: defSSM.Recursive}},
		{"ssm_import:\n  enabled: false", config.SSMImport{Enabled: false, Path: defSSM.Path, Recursive: defSSM.Recursive}},
		{"ssm_import:\n  enabled: true", config.SSMImport{Enabled: true, Path: defSSM.Path, Recursive: defSSM.Recursive}},
		{"ssm_import:\n  enabled: true\n  path: /test/thing", config.SSMImport{Enabled: true, Path: util.Ptr("/test/thing"), Recursive: defSSM.Recursive}},
		{`ssm_import: /path/to/something`, config.SSMImport{Enabled: true, Path: util.Ptr("/path/to/something"), Recursive: &bTrue}},
		{`ssm_import: "/path/to/something/{{ .ProjectName }}"`, config.SSMImport{Enabled: true, Path: util.Ptr("/path/to/something/{{ .ProjectName }}"), Recursive: &bTrue}},
	}

	type conDummy struct {
		SSMImport *config.SSMImport `yaml:"ssm_import,omitempty" json:"ssm_import,omitempty"`
	}

	for _, table := range tables {
		con := conDummy{}
		if err := yaml.UnmarshalStrict([]byte(table.str), &con); err != nil {
			t.Errorf("unexpected error for <%s> %s", table.str, err)
		}

		exp := table.expected
		exp.ApplyDefaults()

		act := con.SSMImport

		if exp.IsEnabled() != act.IsEnabled() {
			t.Errorf("expected <%s> to have IsEnabled=%t but got %t", table.str, exp.IsEnabled(), act.IsEnabled())
		}

		if *exp.Path != *act.Path {
			t.Errorf("expected <%s> to have Path=%v but got %v", table.str, *exp.Path, *act.Path)
		}

		if *exp.Recursive != *act.Recursive {
			t.Errorf("expected <%s> to have Recursive=%v but got %v", table.str, *exp.Recursive, *act.Recursive)
		}
	}
}
