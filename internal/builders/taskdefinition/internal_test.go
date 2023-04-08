package taskdefinition

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil/buildtestutils"
	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
)

func TestContainerTpl(t *testing.T) {
	ctx := buildtestutils.LoadProjectConfig(t, "../testdata/smoke.yml")

	builder, _ := NewBuilder(ctx, ctx.Project.ConsoleTask)

	tables := []struct {
		expected string
		obj      any
	}{
		{"Foo", "Foo"},
		{"Bar", util.Ptr("Bar")},
		{"Baz", ecsTypes.ContainerDefinition{Name: util.Ptr("Baz")}},
		{"Boo", &ecsTypes.ContainerDefinition{Name: util.Ptr("Boo")}},
		{"Yar", &config.CommonContainerAttrs{Name: "Yar"}},
	}

	for _, table := range tables {
		t.Run(table.expected, func(t *testing.T) {
			result, err := builder.containerTpl(table.obj).Apply("{{.Container}}")
			require.NoError(t, err)
			require.Equal(t, table.expected, result)
		})
	}

	t.Run("bad param", func(t *testing.T) {
		require.Panics(t, func() {
			builder.containerTpl(&config.Project{})
		})
	})

}
