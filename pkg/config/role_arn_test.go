package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func TestRoleArn(t *testing.T) {

	closeMock := testutil.MockSimpleStsProxy(t)
	defer closeMock()

	ctx, err := config.NewFromYAML("testdata/simple.yml")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	tables := []struct {
		str  string
		name string
		arn  string
	}{
		{"fakerole", "fakerole", "arn:aws:iam::555555555555:role/fakerole"},
		{"arn:aws:iam::1234567890:role/faker2", "faker2", "arn:aws:iam::1234567890:role/faker2"},
	}

	for _, table := range tables {
		roleArn, err := yaml.ParseYAMLString[config.RoleArn](table.str)
		if err != nil {
			t.Errorf("RoleStr <%s> gave error: %s", table.str, err)
			break
		}

		nameVal, err := roleArn.Name(ctx)
		if err != nil {
			t.Errorf("RoleStr <%s> gave error during name eval: %s", table.str, err)
			break
		}

		if nameVal != table.name {
			t.Errorf("RoleStr <%s> Name Mismatch expected=%s got=%s", table.str, table.name, nameVal)
		}

		arnVal, err := roleArn.Arn(ctx)
		if err != nil {
			t.Errorf("RoleStr <%s> gave error during arn eval: %s", table.str, err)
			break
		}

		if arnVal != table.arn {
			t.Errorf("RoleStr <%s> ARN Mismatch expected=%s got=%s", table.str, table.arn, arnVal)
		}

	}
}
