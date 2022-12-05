package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/stretchr/testify/require"
)

func TestDependsOn_NewFromString_Valid(t *testing.T) {

	st := NewSchemaTester[config.DependsOn](t, &config.DependsOn{})

	tables := []struct {
		str      string
		expected config.DependsOn
	}{
		{`"thing:START"`, config.DependsOn{Condition: ecsTypes.ContainerConditionStart, Name: aws.String("thing")}},
		{`"foo"`, config.DependsOn{Condition: ecsTypes.ContainerConditionStart, Name: aws.String("foo")}},
	}

	for _, table := range tables {
		st.AssertValid(table.str, true)
		obj, err := st.Parse(table.str)

		require.NoError(t, err)
		st.AssertMatchExpected(obj, table.expected, true)
	}
}

func TestDependsOn_NewFromString_Invalid(t *testing.T) {

	st := NewSchemaTester[config.DependsOn](t, &config.DependsOn{})

	tables := []struct {
		str string
	}{
		{`"thing:START:something"`},
		{`"!!! whatever"`},
		{`"whatever thing"`},
	}

	for _, table := range tables {
		valid := st.AssertValid(table.str, false)
		require.Falsef(t, valid, "expected: <%s> to not be valid, but it was", table.str)
		// _, err := st.Parse(table.str)
		// if err == nil {
		// 	t.Errorf("expected <%s> to fail", table.str)
		// }
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

		dep, err := config.NewDependsOnFromString(table.str)

		if table.expectedErr != "" {
			require.Error(t, err)
			require.ErrorContains(t, err, table.expectedErr)
			continue
		}

		require.NoError(t, err)

		require.EqualValues(t, table.deps.Name, dep.Name)
		require.EqualValues(t, table.deps.Condition, dep.Condition)
	}
}
