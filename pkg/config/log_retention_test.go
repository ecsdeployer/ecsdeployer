package config_test

import (
	"testing"

	"ecsdeployer.com/ecsdeployer/internal/util"
	"ecsdeployer.com/ecsdeployer/internal/yaml"
	"ecsdeployer.com/ecsdeployer/pkg/config"
	logTypes "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
	"github.com/stretchr/testify/require"
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

		t.Run(table.str, func(t *testing.T) {

			obj, err := yaml.ParseYAMLString[dummy](table.str)

			if !table.valid {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)

			ret := obj.Retention

			require.Equalf(t, table.forever, ret.Forever(), "Forever")

			if table.forever {
				return
			}
			require.Equalf(t, table.days, ret.Days(), "Days")
			require.EqualValuesf(t, &table.days, ret.ToAwsInt32(), "ToAwsInt32")
		})
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
		require.NoError(t, err)
		st.AssertMatchExpected(obj, table.expected, true)
	}
}

func TestParseLogRetention(t *testing.T) {
	t.Run("normal", func(t *testing.T) {
		tables := []struct {
			val  config.LogRetention
			days int32
		}{
			{util.Must(config.ParseLogRetention("forever")), -1},
			{util.Must(config.ParseLogRetention(int32(180))), 180},
			{util.Must(config.ParseLogRetention(int64(300))), 300},
			{util.Must(config.ParseLogRetention(int(180))), 180},
			{util.Must(config.ParseLogRetention("10")), 10},
		}

		for _, table := range tables {
			require.Equal(t, table.days, table.val.Days())
		}
	})

	t.Run("invalids", func(t *testing.T) {
		tables := []struct {
			genFunc func() (any, error)
		}{
			{func() (any, error) { return config.ParseLogRetention("never") }},
			{func() (any, error) { return config.ParseLogRetention(-1) }},
		}

		for _, table := range tables {
			_, err := table.genFunc()
			require.Error(t, err)
		}
	})

}

func TestLogRetention_EqualsLogGroup(t *testing.T) {

	buildLogGroup := func(days int32) logTypes.LogGroup {
		lg := logTypes.LogGroup{}
		if days >= 0 {
			lg.RetentionInDays = &days
		}
		return lg
	}

	tables := []struct {
		val config.LogRetention
		lg  logTypes.LogGroup
	}{
		{util.Must(config.ParseLogRetention("forever")), buildLogGroup(-1)},
		{util.Must(config.ParseLogRetention(int32(180))), buildLogGroup(180)},
		{util.Must(config.ParseLogRetention(int64(300))), buildLogGroup(300)},
		{util.Must(config.ParseLogRetention(int(180))), buildLogGroup(180)},
		{util.Must(config.ParseLogRetention("10")), buildLogGroup(10)},
	}

	for _, table := range tables {
		require.True(t, table.val.EqualsLogGroup(table.lg))
	}
}
