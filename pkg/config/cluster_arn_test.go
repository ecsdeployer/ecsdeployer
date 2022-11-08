package config_test

import (
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

	for _, table := range tables {
		clusterArn, err := yaml.ParseYAMLString[config.ClusterArn](table.str)
		if err != nil {
			t.Errorf("ClusterStr <%s> gave error: %s", table.str, err)
			break
		}

		nameVal, err := clusterArn.Name(ctx)
		if err != nil {
			t.Errorf("ClusterStr <%s> gave error during name eval: %s", table.str, err)
			break
		}

		if nameVal != table.name {
			t.Errorf("ClusterStr <%s> Name Mismatch expected=%s got=%s", table.str, table.name, nameVal)
		}

		arnVal, err := clusterArn.Arn(ctx)
		if err != nil {
			t.Errorf("ClusterStr <%s> gave error during arn eval: %s", table.str, err)
			break
		}

		if arnVal != table.arn {
			t.Errorf("ClusterStr <%s> ARN Mismatch expected=%s got=%s", table.str, table.arn, arnVal)
		}

	}
}
