package config_test

import (
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestRoleArn(t *testing.T) {

	testutil.MockSimpleStsProxy(t)

	ctx, err := config.NewFromYAML("testdata/simple.yml")
	require.NoError(t, err)

	tables := []struct {
		str  string
		name string
		arn  string
	}{
		{"fakerole", "fakerole", "arn:aws:iam::555555555555:role/fakerole"},
		{"arn:aws:iam::1234567890:role/faker2", "faker2", "arn:aws:iam::1234567890:role/faker2"},
	}

	for testNum, table := range tables {
		t.Run(fmt.Sprintf("test_%02d", testNum+1), func(t *testing.T) {

			roleArn, err := yaml.ParseYAMLString[config.RoleArn](table.str)
			require.NoError(t, err)

			nameVal, err := roleArn.Name(ctx)
			require.NoErrorf(t, err, "Failure during name eval for '%s'", table.str)
			require.Equal(t, table.name, nameVal)

			arnVal, err := roleArn.Arn(ctx)
			require.NoError(t, err)
			require.Equal(t, table.arn, arnVal)
		})

	}
}
