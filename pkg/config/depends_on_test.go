package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	ecsTypes "github.com/aws/aws-sdk-go-v2/service/ecs/types"
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
		if err != nil {
			t.Errorf("error: %s", err)
		}
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
		if valid {
			t.Errorf("expected: <%s> to not be valid, but it was", table.str)
		}
		// _, err := st.Parse(table.str)
		// if err == nil {
		// 	t.Errorf("expected <%s> to fail", table.str)
		// }
	}
}
