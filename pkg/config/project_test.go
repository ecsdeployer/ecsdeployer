package config_test

import (
	"os"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
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
	proj := &config.Project{}
	proj.ApplyDefaults()

	// tables := []struct{}{}

	// for _, table := range tables {

	// }
}

func TestProject_Unmarshal(t *testing.T) {
	// tables := []struct {
	// 	str     string
	// 	enabled bool
	// }{
	// 	{"console: true", true},
	// 	{"console: false", false},
	// 	{"console:\n  enabled: true", true},
	// 	{"console:\n  enabled: false", false},
	// }

	// type conDummy struct {
	// 	Console *config.ConsoleTask `yaml:"console,omitempty" json:"console,omitempty"`
	// }

	// for _, table := range tables {
	// 	con := conDummy{}
	// 	if err := yaml.UnmarshalStrict([]byte(table.str), &con); err != nil {
	// 		t.Errorf("unexpected error for <%s> %s", table.str, err)
	// 	}

	// 	if table.enabled != con.Console.IsEnabled() {
	// 		t.Errorf("expected <%s> to %v console", table.str, table.enabled)
	// 	}
	// }
	t.Skip("Not finished")
}

// make sure that our doc examples are all valid
func TestProject_SchemaCheck_Examples(t *testing.T) {

	st := NewSchemaTester[config.Project](t, config.Project{})
	sc := testutil.NewSchemaChecker(&config.Project{})

	tables := []struct {
		filepath string
	}{
		{"../../cmd/testdata/valid.yml"},
		{"../../www/docs/static/examples/generic.yml"},
		{"../../www/docs/static/examples/simple_web.yml"},
	}

	for _, table := range tables {
		obj, err := yaml.ParseYAMLFile[config.Project](table.filepath)
		require.NoErrorf(t, err, "File %s", table.filepath)

		st.AssertValidObj(*obj, true)

		err = obj.Validate()
		require.NoErrorf(t, err, "File %s failed validation", table.filepath)

		fileData, err := os.ReadFile(table.filepath)
		require.NoError(t, err, "Unable to read test file??")
		require.NoError(t, sc.CheckYAML(t, string(fileData)))

	}
}

// func TestProject_Validate(t *testing.T) {
// 	t.Skip("Not finished")
// }
