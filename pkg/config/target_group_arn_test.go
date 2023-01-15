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

	testutil.StartMocker(t, &awsmocker.MockerOptions{
		Mocks: []*awsmocker.MockedEndpoint{
			testutil.Mock_ELBv2_DescribeTargetGroups_Single_Success("dummytg"),
			testutil.Mock_ELBv2_DescribeTargetGroups_Single_Failure("invalid"),
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
		{"invalid", "invalid", ""},
		// {"arn:aws:iam::1234567890:role/faker2", "faker2", "arn:aws:iam::1234567890:role/faker2"},
	}

	for _, table := range tables {
		t.Run(table.str, func(t *testing.T) {
			tgArn, err := yaml.ParseYAMLString[config.TargetGroupArn](table.str)
			require.NoError(t, err)

			nameVal, err := tgArn.Name(ctx)
			require.NoErrorf(t, err, "name eval")

			require.Equalf(t, table.name, nameVal, "Name Mismatch")

			if table.arn == "" {
				_, err := tgArn.Arn(ctx)
				require.Errorf(t, err, "arn eval")
			} else {
				arnVal, err := tgArn.Arn(ctx)
				require.NoErrorf(t, err, "arn eval")
				require.Equalf(t, table.arn, arnVal, "ARN Mismatch")
			}

			data, err := json.Marshal(tgArn)
			require.NoError(t, err)
			if table.arn != "" {
				require.Equal(t, fmt.Sprintf(`"%s"`, table.arn), string(data))

			}
		})

	}
}
