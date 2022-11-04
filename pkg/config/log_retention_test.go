package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
)

func TestLogRetention_Unmarshal(t *testing.T) {

	tables := []struct {
		str     string
		valid   bool
		forever bool
		days    int32
	}{
		// valid
		{`retention: forever`, true, true, -1},
		{`retention: "forever"`, true, true, -1},
		{`retention: 1`, true, false, 1},
		{`retention: 365`, true, false, 365},

		// Invalid ones
		{`retention: 0`, false, false, -1},
		{`retention: -1`, false, false, 0},
		{`retention: 1.2`, false, false, 0},
		{`retention: false`, false, false, 0},
		{`retention: true`, false, false, 0},
		{`retention: "always"`, false, false, 0},
	}

	type dummy struct {
		Retention *config.LogRetention `yaml:"retention,omitempty" json:"retention,omitempty"`
	}

	for _, table := range tables {
		obj := dummy{}
		err := yaml.UnmarshalStrict([]byte(table.str), &obj)

		if !table.valid {
			if err == nil {
				t.Errorf("Expected <%s> to cause an error, did not get one", table.str)
			}

			continue
		}

		if err != nil {
			t.Errorf("unexpected error for <%s> %s", table.str, err)
		}

		ret := obj.Retention

		if ret.Forever() != table.forever {
			t.Errorf("expected <%s> to have forever=%v but got %v", table.str, table.forever, ret.Forever())
			continue
		}

		if table.forever {
			continue
		}

		if table.days != ret.Days() {
			t.Errorf("expected <%s> to have days=%v but got %v", table.str, table.days, ret.Days())
		}

	}
}

func TestLogRetention_Schema(t *testing.T) {
	st := NewSchemaTester[config.LogRetention](t, &config.LogRetention{})

	tables := []struct {
		str      string
		expected config.LogRetention
	}{
		{`"forever"`, util.Must(config.ParseLogRetention("forever"))},
		{`"123"`, util.Must(config.ParseLogRetention(123))},
		{`123`, util.Must(config.ParseLogRetention(123))},
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
