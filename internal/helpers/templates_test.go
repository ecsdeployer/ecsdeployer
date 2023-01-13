package helpers_test

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/helpers"
	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestGetTemplatedPrefix(t *testing.T) {
	proj, err := config.LoadFromBytes([]byte(testutil.CleanTestYaml(`
	project: myproject
	cluster: mycluster`)))
	require.NoError(t, err, "Failed to load project... bad yaml?")

	ctx := config.New(proj)

	tables := []struct {
		str      string
		expected string
	}{
		{"", ""},
		{"prefix-{{.TaskName}}", "prefix-"},
		{"prefix-{{.TaskName}}-suffix", "prefix-"},
		{"{{.TaskName}}-suffix", ""},
		{"{{.TaskName}}-suffix", ""},
		{"/{{.ProjectName}}/stuff/{{.TaskName}}-suffix", "/myproject/stuff/"},
		{"/test/thing", "/test/thing"},
	}

	for tNum, table := range tables {
		t.Run(fmt.Sprintf("test_%02d", tNum+1), func(t *testing.T) {
			outStr, err := helpers.GetTemplatedPrefix(ctx, table.str)
			require.NoError(t, err)

			require.Equal(t, table.expected, outStr)
		})
	}
}

func TestGetDefaultTaskTemplateFields(t *testing.T) {

	proj, err := config.LoadFromBytes([]byte(testutil.CleanTestYaml(`
	project: test
	cluster: test`)))
	require.NoError(t, err, "Failed to load project... bad yaml?")

	ctx := config.New(proj)

	common := &config.CommonTaskAttrs{
		CommonContainerAttrs: config.CommonContainerAttrs{
			Name: "dummy",
		},
	}

	fields, err := helpers.GetDefaultTaskTemplateFields(ctx, common)
	require.NoError(t, err)

	require.Equal(t, "dummy", fields["Name"])
	require.Equal(t, "dummy", fields["TaskName"])
	require.Equal(t, "amd64", fields["Arch"])
}
