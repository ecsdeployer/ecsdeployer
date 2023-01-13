package config_test

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestProject_ApplyDefaults(t *testing.T) {
	proj := &config.Project{}
	proj.ApplyDefaults()

	require.NotNil(t, proj.ConsoleTask, "ConsoleTask")

	require.NotNil(t, proj.Image, "Image")
	require.Equal(t, "{{ .Image }}", proj.Image.Value())
}

func TestProject_Validate(t *testing.T) {
	tables := []struct {
		str      string
		errMatch string
	}{
		{
			str:      "project: fake\ncluster: fake",
			errMatch: "",
		},
		{
			str:      `cluster: fake`,
			errMatch: "must provide a project name",
		},
		{
			str:      `project: "bad name"`,
			errMatch: "Project name must be lower",
		},
		{
			str: `
			project: fake
			stage: "bad stage"`,
			errMatch: "Stage name must be lower",
		},
		{
			str: `
			project: fake
			cluster: fake
			predeploy:
				- name: thing
				- name: thing
				- name: thing2`,
			errMatch: "Duplicate Resource Name",
		},
		{
			str: `
			project: fake
			cronjobs:
				- name: thing
					schedule: rate(1)`,
			errMatch: "provide a CronLauncher role",
		},
		{
			str:      `project: fake`,
			errMatch: "provide a cluster",
		},
	}

	for tNum, table := range tables {
		t.Run(fmt.Sprintf("test_%02d", tNum+1), func(t *testing.T) {
			cleanStr := testutil.CleanTestYaml(table.str)

			fmt.Println(cleanStr)

			proj, err := config.LoadFromBytes([]byte(cleanStr))
			if table.errMatch == "" {
				require.NoError(t, err)
				require.NotNil(t, proj)
				return
			}

			require.Error(t, err)
			require.ErrorContains(t, err, table.errMatch)
		})
	}
}

func TestProject_Loading(t *testing.T) {

	tables := []struct {
		filepath string
	}{
		{"../../cmd/testdata/valid.yml"},
		{"../../www/docs/static/examples/generic.yml"},
		{"../../www/docs/static/examples/simple_web.yml"},
	}

	for _, table := range tables {
		t.Run(table.filepath, func(t *testing.T) {
			obj, err := config.Load(table.filepath)
			require.NoError(t, err)

			fileData, err := os.ReadFile(table.filepath)
			require.NoError(t, err)
			strReader := strings.NewReader(string(fileData))

			obj2, err := config.LoadReader(strReader)
			require.NoError(t, err)

			json1, _ := util.Jsonify(obj)
			json2, _ := util.Jsonify(obj2)

			require.JSONEq(t, json1, json2)

		})
	}
}

// make sure that our doc examples are all valid
func TestProject_SchemaCheck_Examples(t *testing.T) {

	sc := testutil.NewSchemaChecker(&config.Project{})

	tables := []struct {
		filepath string
	}{
		{"../../cmd/testdata/valid.yml"},
		{"../../www/docs/static/examples/generic.yml"},
		{"../../www/docs/static/examples/simple_web.yml"},
	}

	for _, table := range tables {
		t.Run(table.filepath, func(t *testing.T) {

			obj, err := yaml.ParseYAMLFile[config.Project](table.filepath)
			require.NoErrorf(t, err, "File %s", table.filepath)

			err = obj.Validate()
			require.NoErrorf(t, err, "File %s failed validation", table.filepath)

			fileData, err := os.ReadFile(table.filepath)
			require.NoError(t, err, "Unable to read test file??")
			require.NoError(t, sc.CheckYAML(t, string(fileData)))
		})
	}
}
