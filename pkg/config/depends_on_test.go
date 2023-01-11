package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/testutil"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
)

func TestDependsOn_Parsing(t *testing.T) {
	sc := testutil.NewSchemaChecker(&config.DependsOn{})

	tables := []struct {
		str      string
		expected *config.DependsOn
	}{
		{`"thing:START"`, &config.DependsOn{Condition: ecsTypes.ContainerConditionStart, Name: aws.String("thing")}},
		{`"foo"`, &config.DependsOn{Condition: ecsTypes.ContainerConditionStart, Name: aws.String("foo")}},

		{`"thing:START:something"`, nil},
		{`"!!! whatever"`, nil},
		{`"whatever thing"`, nil},
	}

	for _, table := range tables {
		t.Run(table.str, func(t *testing.T) {
			obj, err := yaml.ParseYAMLString[config.DependsOn](table.str)

			if table.expected == nil {
				// this is a bit lazy, but it means we don't have to make the validation or schema be super picky.
				// we can just rely on aws rejecting them
				require.Truef(t, (err != nil) || sc.CheckYAML(t, table.str) != nil, "Should have either failed parsing or a schema validation")
				return
			}

			require.NoError(t, err)
			require.NoError(t, sc.CheckYAML(t, table.str))
			require.NotNil(t, obj)

			require.Equal(t, table.expected.Condition, obj.Condition)
			require.Equal(t, *table.expected.Name, *obj.Name)

		})
	}
}

func TestNewDependsOnFromString(t *testing.T) {
	tables := []struct {
		str         string
		expectedErr string
		deps        *config.DependsOn
	}{
		{"thing:START", "", &config.DependsOn{ecsTypes.ContainerConditionStart, aws.String("thing")}},
		{"foo", "", &config.DependsOn{ecsTypes.ContainerConditionStart, aws.String("foo")}},

		// failures
		{"thing:START:something", "must be object, or string of", nil},
		{"thing:STORT", "not a valid condition", nil},
	}

	for _, table := range tables {
		t.Run(table.str, func(t *testing.T) {

			dep, err := config.NewDependsOnFromString(table.str)

			if table.expectedErr != "" {
				require.Error(t, err)
				require.ErrorContains(t, err, table.expectedErr)
				require.ErrorIs(t, err, config.ErrValidation)
				return
			}

			require.NoError(t, err)

			require.EqualValues(t, table.deps.Name, dep.Name)
			require.EqualValues(t, table.deps.Condition, dep.Condition)
		})
	}
}
