package config_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	"github.com/stretchr/testify/require"
)

func TestEnvVars(t *testing.T) {

	someValue := "something"
	val := &someValue

	tables := []struct {
		isTpl      bool
		isSSM      bool
		isPlain    bool
		isUnset    bool
		marshalStr string
		obj        *config.EnvVar
	}{
		// {false, false, false, false, `""`, &config.EnvVar{}},
		{true, false, false, false, `{"template":"something"}`, &config.EnvVar{ValueTemplate: val}},
		{false, true, false, false, `{"ssm":"something"}`, &config.EnvVar{ValueSSM: val}},
		{false, false, true, false, `{"value":"something"}`, &config.EnvVar{Value: val}},
		{false, false, false, true, `{"unset":true}`, &config.EnvVar{Unset: true}},
	}
	for i, table := range tables {
		t.Run(fmt.Sprintf("entry_%02d", i), func(t *testing.T) {
			require.Equalf(t, table.isTpl, table.obj.IsTemplated(), "IsTemplated")
			require.Equalf(t, table.isSSM, table.obj.IsSSM(), "IsSSM")
			require.Equalf(t, table.isPlain, table.obj.IsPlain(), "IsPlain")
			require.Equalf(t, table.isUnset, table.obj.IsUnset(), "IsUnset")

			require.NoErrorf(t, table.obj.Validate(), "Validate")

			jsonData, err := json.Marshal(table.obj)
			require.NoError(t, err)

			require.JSONEq(t, table.marshalStr, string(jsonData))

		})
	}

}

func TestEnvVar_UnmarshalYAML(t *testing.T) {

	type evtype int

	const (
		evPlain evtype = iota
		evSSM
		evTemplated
		evUnset
	)

	tables := []struct {
		str      string
		vartype  evtype
		expVal   string
		invalid  bool
		errMatch string
	}{
		{`""`, evPlain, "", false, ""},
		{`testing`, evPlain, "testing", false, ""},
		{`ssm: testing`, evSSM, "testing", false, ""},
		{`template: testing`, evTemplated, "testing", false, ""},
		{`value: testing`, evPlain, "testing", false, ""},
		{`unset: true`, evUnset, "", false, ""},

		{`unset: false`, -1, "", true, "You need to provide some value for an env var"},
		{"ssm: testing\nvalue: thing", -1, "", true, "You can only provide one type"},
	}

	for tNum, table := range tables {
		t.Run(fmt.Sprintf("test_%02d", tNum+1), func(t *testing.T) {

			obj, err := yaml.ParseYAMLString[config.EnvVar](table.str)

			if table.invalid {
				require.Error(t, err)
				if table.errMatch != "" {
					require.ErrorContains(t, err, table.errMatch)
				}
				return
			}

			require.NoError(t, err)
			require.NoError(t, obj.Validate())

			switch table.vartype {
			case evPlain:
				require.True(t, obj.IsPlain())
				require.Equal(t, table.expVal, *obj.Value)
			case evSSM:
				require.True(t, obj.IsSSM())
				require.Equal(t, table.expVal, *obj.ValueSSM)
			case evTemplated:
				require.True(t, obj.IsTemplated())
				require.Equal(t, table.expVal, *obj.ValueTemplate)
			case evUnset:
				require.True(t, obj.IsUnset())
			}

		})
	}
}
