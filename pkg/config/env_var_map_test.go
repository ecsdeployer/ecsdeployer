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

func TestEnvVarMap(t *testing.T) {
	t.Run("Filter", func(t *testing.T) {
		map2 := config.EnvVarMap{
			"THING1": config.NewEnvVar(config.EnvVarTypePlain, "foo"),
			"THING4": config.NewEnvVar(config.EnvVarTypeUnset, ""),
			"THING5": config.NewEnvVar(config.EnvVarTypePlain, "bar"),
			"THING6": config.NewEnvVar(config.EnvVarTypePlain, "test6"),
			"THING7": config.NewEnvVar(config.EnvVarTypePlain, "test7"),
		}

		filtered := map2.Filter()
		require.Len(t, map2, 5)
		require.Len(t, filtered, 4)
		require.Contains(t, map2, "THING4")
		require.NotContains(t, filtered, "THING4")
	})
	t.Run("HasSSM", func(t *testing.T) {
		map1 := config.EnvVarMap{
			"THING1": config.NewEnvVar(config.EnvVarTypePlain, "test1"),
			"THING2": config.NewEnvVar(config.EnvVarTypeSSM, "test2"),
			"THING3": config.NewEnvVar(config.EnvVarTypeTemplated, "tplthing"),
			"THING4": config.NewEnvVar(config.EnvVarTypePlain, "test4"),
			"THING5": config.NewEnvVar(config.EnvVarTypePlain, "test5"),
		}

		map2 := config.EnvVarMap{
			"THING1": config.NewEnvVar(config.EnvVarTypePlain, "foo"),
			"THING4": config.NewEnvVar(config.EnvVarTypeUnset, ""),
			"THING5": config.NewEnvVar(config.EnvVarTypePlain, "bar"),
			"THING6": config.NewEnvVar(config.EnvVarTypePlain, "test6"),
			"THING7": config.NewEnvVar(config.EnvVarTypePlain, "test7"),
		}

		require.True(t, map1.HasSSM())
		require.False(t, map2.HasSSM())
	})
}

func TestMergeEnvVarMaps(t *testing.T) {
	map1 := config.EnvVarMap{
		"THING1": config.NewEnvVar(config.EnvVarTypePlain, "test1"),
		"THING2": config.NewEnvVar(config.EnvVarTypeSSM, "test2"),
		"THING3": config.NewEnvVar(config.EnvVarTypeTemplated, "tplthing"),
		"THING4": config.NewEnvVar(config.EnvVarTypePlain, "test4"),
		"THING5": config.NewEnvVar(config.EnvVarTypePlain, "test5"),
	}

	map2 := config.EnvVarMap{
		"THING1": config.NewEnvVar(config.EnvVarTypePlain, "foo"),
		"THING4": config.NewEnvVar(config.EnvVarTypeUnset, ""),
		"THING5": config.NewEnvVar(config.EnvVarTypePlain, "bar"),
		"THING6": config.NewEnvVar(config.EnvVarTypePlain, "test6"),
		"THING7": config.NewEnvVar(config.EnvVarTypePlain, "test7"),
	}

	map3 := config.EnvVarMap{
		"THING7": config.NewEnvVar(config.EnvVarTypePlain, "newval"),
	}

	newmap := config.MergeEnvVarMaps(map1, map2, map3)

	require.Equal(t, "foo", util.Must(newmap["THING1"].GetValue(testutil.TplDummy)))
	require.Equal(t, "test2", util.Must(newmap["THING2"].GetValue(testutil.TplDummy)))
	require.Equal(t, "tplthing", util.Must(newmap["THING3"].GetValue(testutil.TplDummy)))
	require.True(t, newmap["THING4"].IsUnset())
	require.Equal(t, "bar", util.Must(newmap["THING5"].GetValue(testutil.TplDummy)))
	require.Equal(t, "test6", util.Must(newmap["THING6"].GetValue(testutil.TplDummy)))
	require.Equal(t, "newval", util.Must(newmap["THING7"].GetValue(testutil.TplDummy)))

}

func TestEnvVarMap_UnmarshalYAML(t *testing.T) {
	tables := []struct {
		str        string
		invalid    bool
		errMatch   string
		badSchema  bool
		hasSSM     bool
		extraCheck func(*testing.T, *config.EnvVarMap)
	}{
		{
			str: `SOMEVAR: testing`,
		},
		{
			str:    `SOMETHING: {ssm: /var/whatever}`,
			hasSSM: true,
		},
		{
			str:       `123SOMETHING: badvar`,
			badSchema: true,
		},
		{
			str: `
			THING1: test
			THING2:
				ssm: /blah
			THING3:
				template: "{{.ProjectName}} thing"
			THING4:
				unset: true
			`,
			hasSSM: true,
			extraCheck: func(t *testing.T, evm *config.EnvVarMap) {
				require.Len(t, *evm, 4)
			},
		},
	}

	sc := testutil.NewSchemaChecker(&config.EnvVarMap{})

	for tNum, table := range tables {
		t.Run(fmt.Sprintf("test_%02d", tNum+1), func(t *testing.T) {
			cleanStr := testutil.CleanTestYaml(table.str)

			schemaErr := sc.CheckYAML(t, cleanStr)
			if table.badSchema {
				require.Error(t, schemaErr)
			} else {
				require.NoError(t, schemaErr)
			}

			obj, err := yaml.ParseYAMLString[config.EnvVarMap](cleanStr)
			if table.invalid || table.errMatch != "" {
				require.Error(t, err)

				if table.errMatch != "" {
					require.ErrorContains(t, err, table.errMatch)
				}
				return
			}

			require.NoError(t, err)
			require.IsType(t, &config.EnvVarMap{}, obj)

			require.Equal(t, table.hasSSM, obj.HasSSM(), "HasSSM")

			if table.extraCheck != nil {
				table.extraCheck(t, obj)
			}

		})
	}
}
