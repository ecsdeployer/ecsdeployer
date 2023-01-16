package config_test

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestClusterArn(t *testing.T) {

	testutil.MockSimpleStsProxy(t)

	ctx, err := config.NewFromYAML("testdata/simple.yml")
	require.NoError(t, err)

	tables := []struct {
		str  string
		name string
		arn  string
	}{
		{"fakecluster", "fakecluster", "arn:aws:ecs:us-east-1:555555555555:cluster/fakecluster"},
		{"arn:aws:ecs:us-east-1:1234567890:cluster/faker2", "faker2", "arn:aws:ecs:us-east-1:1234567890:cluster/faker2"},
	}

	for testNum, table := range tables {

		t.Run(fmt.Sprintf("test_%02d", testNum+1), func(t *testing.T) {

			obj, err := yaml.ParseYAMLString[config.ClusterArn](table.str)
			require.NoError(t, err)

			nameVal, err := obj.Name(ctx)
			require.NoErrorf(t, err, "Failure during name eval for '%s'", table.str)
			require.Equal(t, table.name, nameVal)

			arnVal, err := obj.Arn(ctx)
			require.NoError(t, err)
			require.Equal(t, table.arn, arnVal)
		})
	}
}
