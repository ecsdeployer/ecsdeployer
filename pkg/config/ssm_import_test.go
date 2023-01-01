package config_test

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestSSMImport_Unmarshal(t *testing.T) {

	defSSM := config.SSMImport{}
	defSSM.ApplyDefaults()

	bTrue := true
	// bFalse := false

	tables := []struct {
		str      string
		expected *config.SSMImport
	}{
		{"true", &config.SSMImport{Enabled: true, Path: defSSM.Path, Recursive: defSSM.Recursive}},
		{"false", &config.SSMImport{Enabled: false, Path: defSSM.Path, Recursive: defSSM.Recursive}},
		{"enabled: false", &config.SSMImport{Enabled: false, Path: defSSM.Path, Recursive: defSSM.Recursive}},
		{"enabled: true", &config.SSMImport{Enabled: true, Path: defSSM.Path, Recursive: defSSM.Recursive}},
		{"enabled: true\npath: /test/thing", &config.SSMImport{Enabled: true, Path: util.Ptr("/test/thing"), Recursive: defSSM.Recursive}},
		{`/path/to/something`, &config.SSMImport{Enabled: true, Path: util.Ptr("/path/to/something"), Recursive: &bTrue}},
		{`"/path/to/something/{{ .ProjectName }}"`, &config.SSMImport{Enabled: true, Path: util.Ptr("/path/to/something/{{ .ProjectName }}"), Recursive: &bTrue}},

		{"1234", nil},      // interpreted as a string, but must start with slash
		{"test/path", nil}, // ssm path must start with slash
	}

	sc := testutil.NewSchemaChecker(&config.SSMImport{})

	for testNum, table := range tables {
		t.Run(fmt.Sprintf("test_%02d", testNum+1), func(t *testing.T) {
			act, err := yaml.ParseYAMLString[config.SSMImport](table.str)

			if table.expected == nil {
				require.Error(t, sc.CheckYAML(t, table.str))
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			// ensure it passes schema check
			require.NoError(t, sc.CheckYAML(t, table.str))

			exp := table.expected
			exp.ApplyDefaults()

			require.Equalf(t, exp.IsEnabled(), act.IsEnabled(), "IsEnabled")
			require.Equalf(t, *exp.Path, *act.Path, "Path")
			require.Equalf(t, *exp.Recursive, *act.Recursive, "Recursive")
		})
	}
}
