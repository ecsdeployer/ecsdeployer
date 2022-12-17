package config_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
	"github.com/webdestroya/awsmocker"
)

func TestTargetGroupArn_Smoke(t *testing.T) {

	awsmocker.Start(t, &awsmocker.MockerOptions{
		Mocks: []*awsmocker.MockedEndpoint{
			testutil.Mock_ELBv2_DescribeTargetGroups_Single_Success("dummytg"),
		},
	})

	ctx, err := config.NewFromYAML("testdata/simple.yml")
	require.NoError(t, err)

	tables := []struct {
		str  string
		name string
		arn  string
	}{
		{"dummytg", "dummytg", "arn:aws:elasticloadbalancing:us-east-1:555555555555:targetgroup/dummytg/73e2d6bc24d8a067"},
		// {"arn:aws:iam::1234567890:role/faker2", "faker2", "arn:aws:iam::1234567890:role/faker2"},
	}

	for _, table := range tables {
		t.Run(table.str, func(t *testing.T) {
			roleArn, err := yaml.ParseYAMLString[config.TargetGroupArn](table.str)
			require.NoError(t, err)

			nameVal, err := roleArn.Name(ctx)
			require.NoErrorf(t, err, "name eval")

			require.Equalf(t, table.name, nameVal, "Name Mismatch")

			arnVal, err := roleArn.Arn(ctx)
			require.NoErrorf(t, err, "arn eval")

			require.Equalf(t, table.arn, arnVal, "ARN Mismatch")

			data, err := json.Marshal(roleArn)
			require.NoError(t, err)
			require.Equal(t, fmt.Sprintf(`"%s"`, table.arn), string(data))
		})

	}
}
